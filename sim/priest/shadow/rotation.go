package shadow

import (
	//"fmt"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var mbIdb = 0
var dpIdx = 1
var vtIdx = 2
var swdIdx = 3

func (spriest *ShadowPriest) OnGCDReady(sim *core.Simulation) {
	spriest.tryUseGCD(sim)
}

func (spriest *ShadowPriest) OnManaTick(sim *core.Simulation) {
	if spriest.FinishedWaitingForManaAndGCDReady(sim) {
		spriest.tryUseGCD(sim)
	}
}
func (spriest *ShadowPriest) tryUseGCD(sim *core.Simulation) {
	// TODO: probably do something different instead of making it global?
	// some global variables used throughout the code
	var currentWait time.Duration
	var mbDamage float64
	var dpDamage float64
	var vtDamage float64
	var swdDamage float64
	var mfDamage float64
	var overwriteDPS float64

	var currDotTickSpeed float64

	//var remain_fight float64

	//if sim.CurrentTime == 0 && spriest.rotation.PrecastVt {
	//spriest.SpendMana(sim, spriest.VampiricTouch.DefaultCast.Cost, spriest.VampiricTouch.ResourceMetrics)
	//spriest.VampiricTouch.SkipCastAndApplyEffects(sim, spriest.CurrentTarget)
	//}

	// initialize function specific variables
	var spell *core.Spell
	var swStacks float64
	var numswptickstime float64
	var cdDpso float64
	var cdDps float64
	var chosenMfs int
	var num_DP_ticks float64
	var wait1 time.Duration
	var wait2 time.Duration
	var wait time.Duration
	var wait3 time.Duration

	// initialize helpful variables for calculations later
	vtCastTime := spriest.ApplyCastSpeed(time.Millisecond * 1500)
	gcd := spriest.SpellGCD()
	tickLength := spriest.MindFlayTickDuration()

	dotTickSpeed := float64(spriest.ApplyCastSpeed(time.Second * 3))
	remain_fight := float64(sim.GetRemainingDuration())
	castMf2 := 0 // if SW stacks = 3, and we want to get SWP up at 5 stacks exactly, then we want to hard code a MF2
	bestIdx := -1

	// How much time until lust is used?
	deltaTimeBL := core.MaxDuration(0, spriest.BLUsedAt-sim.CurrentTime)
	// How many VT ticks before lust is used?
	numVTbeforeBL := math.Floor(deltaTimeBL.Seconds() / (dotTickSpeed * 1e-9))
	if numVTbeforeBL < 0 {
		numVTbeforeBL = 0
	}

	// Decide if precast MB or VT is more dps
	if sim.CurrentTime == 0 && spriest.rotation.PrecastVt && spriest.CurrentMana() == spriest.MaxMana() {
		if deltaTimeBL >= gcd && numVTbeforeBL < 1 && sim.CurrentTime.Seconds() < float64(spriest.BLUsedAt) {
			spriest.SpendMana(sim, spriest.MindBlast.DefaultCast.Cost, spriest.MindBlast.ResourceMetrics)
			spriest.MindBlast.SkipCastAndApplyEffects(sim, spriest.CurrentTarget)
			spriest.MindBlast.CD.UsePrePull(sim, sim.CurrentTime)
		} else {
			spriest.SpendMana(sim, spriest.VampiricTouch.DefaultCast.Cost, spriest.VampiricTouch.ResourceMetrics)
			spriest.VampiricTouch.SkipCastAndApplyEffects(sim, spriest.CurrentTarget)
		}
	}

	// grab all of the shadow priest spell CDs remaining durations to use in the dps calculation
	allCDs := []time.Duration{
		core.MaxDuration(0, spriest.MindBlast.TimeToReady(sim)),
		core.MaxDuration(0, spriest.DevouringPlagueDot.RemainingDuration(sim)),
		core.MaxDuration(0, spriest.VampiricTouchDot.RemainingDuration(sim)-vtCastTime),
		core.MaxDuration(0, spriest.ShadowWordDeath.TimeToReady(sim)),
		0,
	}

	rotType := spriest.rotation.RotationType

	if spriest.ShadowWeavingAura.IsActive() {
		swStacks = float64(spriest.ShadowWeavingAura.GetStacks())
	}

	if rotType == proto.ShadowPriest_Rotation_AoE {
		numTicks := 5
		spell = spriest.MindSear[numTicks]
		if success := spell.Cast(sim, spriest.CurrentTarget); !success {
			spriest.WaitForMana(sim, spell.CurCast.Cost)
		}
		return

	} else if rotType == proto.ShadowPriest_Rotation_Basic || rotType == proto.ShadowPriest_Rotation_Clipping {

		if spriest.DevouringPlagueDot.RemainingDuration(sim) <= 0 {
			bestIdx = 1
		} else if spriest.Talents.VampiricTouch && spriest.VampiricTouchDot.RemainingDuration(sim) <= vtCastTime {
			bestIdx = 2
		} else if !spriest.ShadowWordPainDot.IsActive() && swStacks >= 5 {
			bestIdx = 5
		} else if spriest.Talents.MindFlay {
			bestIdx = 4
		}
	} else {
		// Reduce number of DP/VT ticks based on remaining duration
		num_DP_ticks = math.Floor(remain_fight / dotTickSpeed)
		if num_DP_ticks > 8 {
			num_DP_ticks = 8
		}

		var blDuration time.Duration
		aura := spriest.GetActiveAuraWithTag(core.BloodlustAuraTag)
		if aura != nil {
			blDuration = aura.RemainingDuration(sim)
		}

		// Spell damage numbers that are updated before each cast in order to determine the most optimal next cast based on dps over a finite window
		// This is needed throughout the code to determine the optimal spell(s) to cast next
		// MB dmg
		mbDamage = 0
		if spriest.options.UseMindBlast {
			mbDamage = spriest.MindBlast.ExpectedDamage(sim, spriest.CurrentTarget)
		}

		// DP dmg
		dpTickDamage := spriest.DevouringPlague.ExpectedDamage(sim, spriest.CurrentTarget)
		dpInitDamage := dpTickDamage * spriest.DpInitMultiplier
		dpDamage = dpInitDamage + dpTickDamage*num_DP_ticks

		// Determine number of DP ticks before BL. If there is at least 1 then it's worth using
		numDPbeforeBL := math.Floor(deltaTimeBL.Seconds() / (dotTickSpeed * 1e-9))
		if (deltaTimeBL > gcd && numDPbeforeBL < 1 && sim.CurrentTime < spriest.BLUsedAt) || (deltaTimeBL <= gcd && deltaTimeBL > time.Millisecond*10) {
			dpDamage = 0
		}

		// VT dmg
		// If there is at least 2 VT ticks then it's worth using
		vtDamage = 0
		if deltaTimeBL <= gcd || numVTbeforeBL >= 2 || sim.CurrentTime >= spriest.BLUsedAt {
			vtDamage = spriest.VampiricTouch.ExpectedDamage(sim, spriest.CurrentTarget)
		}

		// SWD dmg
		swdDamage = 0
		if spriest.options.UseShadowWordDeath {
			swdDamage = spriest.ShadowWordDeath.ExpectedDamage(sim, spriest.CurrentTarget)
		}

		mfDamage = spriest.MindFlay[3].ExpectedDamage(sim, spriest.CurrentTarget)
		swpTickDamage := spriest.ShadowWordPain.ExpectedDamage(sim, spriest.CurrentTarget)

		// this should be cleaned up, but essentially we want to cast SWP either 3rd or 5th in the rotation which is fight length dependent
		castSwpNow := false // if SW stacks = 3, and we want to get SWP up immediately becaues fight length is low enough, then this flag gets set to 1
		if swStacks > 2 && swStacks < 5 && !spriest.ShadowWordPainDot.IsActive() {
			addedDmg := mbDamage*0.12 + mfDamage*0.22*2/3 + swpTickDamage*2*gcd.Seconds()/3
			numswptickstime = addedDmg / (swpTickDamage * 0.06) * 3 //if the fight lenght is < numswptickstime then use swp 3rd.. if > then use at weaving = 5
			//
			if remain_fight*math.Pow(10, -9) < numswptickstime { //
				castSwpNow = true
			} else {
				castMf2 = 1
			}
		}

		var currDPS float64
		var nextTickWait time.Duration
		var currDPS2 float64
		var overwriteDPS2 float64

		if spriest.DevouringPlagueDot.IsActive() {
			nextTickWait = spriest.DevouringPlagueDot.TimeUntilNextTick(sim)

			dpDotCurr := spriest.DevouringPlague.ExpectedDamageFromCurrentSnapshot(sim, spriest.CurrentTarget)
			dpInitCurr := dpTickDamage * spriest.DpInitMultiplier

			cdDamage := mbDamage
			if spriest.T10FourSetBonus || cdDamage == 0 {
				cdDamage = mfDamage / 3 * 2
			}

			currDotTickSpeed = spriest.DevouringPlagueDot.TickPeriod().Seconds()
			dotTickSpeednew := 3 / spriest.CastSpeed
			currDPS = (dpInitCurr + dpDotCurr*8 + cdDamage) / (currDotTickSpeed * 8)
			overwriteDPS = (dpInitCurr + dpInitDamage + dpDotCurr*1 + dpTickDamage) / (dotTickSpeednew*8 + currDotTickSpeed*1)

			if blDuration.Seconds() < 3 && blDuration.Seconds() > 0.1 {
				dpRemainTicks := 8 - allCDs[dpIdx].Seconds()/currDotTickSpeed
				overwriteDPS2 = dpInitCurr + dpRemainTicks*(dpDotCurr-dpDotCurr/spriest.CastSpeed)
				currDPS2 = cdDamage

				// if sim.Log != nil {
				// 	spriest.Log(sim, "currDPS2[%d]", currDPS2)
				// 	spriest.Log(sim, "overwriteDPS2[%d]", overwriteDPS2)
				// 	spriest.Log(sim, "dpRemainTicks[%d]", dpRemainTicks)
				// }
			}
		}

		// Make an array of DPCT per spell that will be used to find the optimal spell to cast
		spellDPCT := []float64{
			// MB dps
			mbDamage / (gcd + allCDs[mbIdb]).Seconds(),
			// DP dps
			dpDamage / (gcd + allCDs[dpIdx]).Seconds(),
			// VT dps
			vtDamage / (gcd + allCDs[vtIdx]).Seconds(),
			// SWD dps
			swdDamage / (gcd + allCDs[swdIdx]).Seconds(),
			// MF dps 3 ticks
			mfDamage / (tickLength * 3).Seconds(),
		}

		//if sim.Log != nil {
		//spriest.Log(sim, "mbDamage[%d]", mbDamage)
		//spriest.Log(sim, "mb time[%d]", allCDs[mbIdb])
		//spriest.Log(sim, "mftime[%d]", float64((tickLength * 3).Seconds()))
		//spriest.Log(sim, "gcd[%d]", gcd.Seconds())
		//spriest.Log(sim, "CastSpeedMultiplier[%d]", spriest.PseudoStats.CastSpeedMultiplier)
		//spriest.Log(sim, "critChance[%d]", critChance)
		//}

		// Find the maximum DPCT spell
		bestDmg := 0.0
		for i, v := range spellDPCT {
			if sim.Log != nil {
				spriest.Log(sim, "\tspellDPCT[%d]: %01.f", i, v)
				spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, allCDs[i].Seconds())
			}
			if v > bestDmg {
				bestIdx = i
				bestDmg = v
			}
		}
		// Find the minimum CD ability to make sure that shouldnt be cast first
		nextCD := core.NeverExpires
		nextIdx := -1
		for i, v := range allCDs[1 : len(allCDs)-1] {
			//	if sim.Log != nil {
			// spriest.Log(sim, "\tallCDs[%d]: %01.f", i, v)
			// 	 spriest.Log(sim, "\tcdDiffs[%d]: %0.1f", i, cdDiffs[i].Seconds())
			//}
			if v < nextCD {
				nextCD = v
				nextIdx = i + 1
			}
		}
		waitmin := nextCD

		// Now it's possible that the wait time for the chosen spell is long, if that's the case, then it might be better to investigate the dps over a 2 spell window to see if casting something else will benefit
		if bestIdx < 4 {
			currentWait = allCDs[bestIdx]
		}

		if allCDs[0] < gcd && bestIdx == 4 && allCDs[3] == 0 {
			totalDps__poss := (mbDamage + swdDamage) / (gcd + gcd).Seconds()
			totalDps__poss3 := (mbDamage + mfDamage*2/3) / (2*tickLength + gcd).Seconds()

			if totalDps__poss > totalDps__poss3 {
				bestIdx = 3
				currentWait = allCDs[bestIdx]
			}
		}

		if nextIdx != 4 && bestIdx != 4 && bestIdx != 5 && currentWait > waitmin && currentWait.Seconds() < 3 { // right now 3 might not be correct number, but we can study this to optimize

			if bestIdx == 2 { // MB VT DP SWD
				cdDpso = vtDamage / (gcd + currentWait).Seconds()
			} else if bestIdx == 0 {
				cdDpso = mbDamage / (gcd + currentWait).Seconds()
			} else if bestIdx == 3 {
				cdDpso = swdDamage / (gcd + currentWait).Seconds()
			} else if bestIdx == 1 {
				cdDpso = dpDamage / (gcd + currentWait).Seconds()
			}

			if nextIdx == 2 {
				cdDps = vtDamage / (gcd + waitmin).Seconds()
			} else if nextIdx == 0 {
				cdDps = mbDamage / (gcd + waitmin).Seconds()
			} else if nextIdx == 3 {
				cdDps = swdDamage / (gcd + waitmin).Seconds()
			} else if nextIdx == 1 {
				cdDps = dpDamage / (gcd + waitmin).Seconds()
			}

			residualWait := currentWait - gcd
			if residualWait < 0 {
				residualWait = 0
			}
			totalDps__poss0 := (cdDpso * (currentWait + gcd).Seconds()) / (gcd + currentWait).Seconds()
			totalDps__poss1 := (cdDpso*(currentWait+gcd).Seconds() + cdDps*(waitmin+gcd).Seconds()) / (waitmin + gcd + gcd + residualWait).Seconds()

			totalDps__poss2 := float64(0)
			totalDps__poss3 := float64(0)

			residualMF := currentWait - 2*tickLength
			if residualMF < 0 {
				residualMF = 0
			}
			totalDps__poss2 = (cdDpso*(currentWait+gcd).Seconds() + mfDamage*2/3) / (2*tickLength + gcd + residualMF).Seconds()
			residualMF = currentWait - 3*tickLength
			if residualMF < 0 {
				residualMF = 0
			}
			totalDps__poss3 = (cdDpso*(currentWait+gcd).Seconds() + mfDamage) / (3*tickLength + gcd + residualMF).Seconds()

			//if sim.Log != nil {
			//spriest.Log(sim, "nextIdx[%d]", nextIdx)
			//spriest.Log(sim, "bestIdx[%d]", bestIdx)
			//spriest.Log(sim, "currentWait[%d]", currentWait.Seconds())
			//spriest.Log(sim, "total_dps__poss0[%d]", totalDps__poss0)
			//spriest.Log(sim, "total_dps__poss1[%d]", totalDps__poss1)
			//spriest.Log(sim, "total_dps__poss2[%d]", totalDps__poss2)
			//spriest.Log(sim, "total_dps__poss3[%d]", totalDps__poss3)
			//}

			// TODO looks fishy, repeated bestIdx = 4
			if (totalDps__poss1 > totalDps__poss0) || (totalDps__poss2 > totalDps__poss0) || (totalDps__poss3 > totalDps__poss0) {
				if totalDps__poss1 > totalDps__poss0 && totalDps__poss1 > totalDps__poss2 && totalDps__poss1 > totalDps__poss3 {
					bestIdx = nextIdx // if choosing the minimum wait time spell first is highest dps, then change the index and current wait
					currentWait = waitmin
				} else if totalDps__poss2 > totalDps__poss0 && totalDps__poss2 > totalDps__poss1 && totalDps__poss2 > totalDps__poss3 {
					//bestIdx = bestIdx // if choosing the minimum wait time spell first is highest dps, then change the index and current wait
					//currentWait = currentWait
					bestIdx = 4
				} else if totalDps__poss3 > totalDps__poss0 && totalDps__poss3 > totalDps__poss1 && totalDps__poss3 > totalDps__poss2 {
					//bestIdx = bestIdx // if choosing the minimum wait time spell first is highest dps, then change the index and current wait
					//currentWait = currentWait
					bestIdx = 4
				} else {
					bestIdx = 4
				}
			}

		}

		// If VT isnt chosen, and reapplying DP is more dps, then overwrite it next
		if overwriteDPS-currDPS > 200 && bestIdx != 2 {
			bestIdx = 1
			currentWait = nextTickWait
		} else {
			overwriteDPS = 0
		}

		// Now it's possible that the wait time is > 1 gcd and is the minimum wait time. this is unlikely in wrath given how good MF is, but still might be worth to check
		// if chosen wait time is > 0.3*GCD (this was optimized in private sim, but might want to reoptimize with procs/ect) then check if it's more dps to add a mf sequence
		if bestIdx != 4 && currentWait.Seconds() > 0.3*gcd.Seconds() {

			if bestIdx == 2 { // MB VT DP SWD
				cdDpso = vtDamage
			} else if bestIdx == 0 {
				cdDpso = mbDamage
			} else if bestIdx == 3 {
				cdDpso = swdDamage
			} else if bestIdx == 1 {
				cdDpso = dpDamage
			}

			addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
			addedgcdtime := addedgcd - time.Duration(2)*tickLength

			deltaMf1 := currentWait - gcd
			if deltaMf1 < 0 {
				deltaMf1 = 0
			}
			deltaMf2 := currentWait - (tickLength*2 + addedgcdtime)
			if deltaMf2 < 0 {
				deltaMf2 = 0
			}
			deltaMf3 := currentWait - tickLength*3
			if deltaMf3 < 0 {
				deltaMf3 = 0
			}

			dpsPossibleshort := []float64{
				(cdDpso) / (gcd + currentWait).Seconds(),
				(cdDpso + mfDamage/3) / (deltaMf1 + gcd + gcd).Seconds(),
				(cdDpso + mfDamage/3*2) / (deltaMf2 + tickLength*2 + addedgcdtime + gcd).Seconds(),
				(cdDpso + mfDamage) / (deltaMf3 + tickLength*3 + gcd).Seconds(),
			}

			// Find the highest possible dps and its index
			highestPossibleIdx := 0
			highestPossibleDmg := 0.0
			for i, v := range dpsPossibleshort {
				if v >= highestPossibleDmg {
					//if sim.Log != nil {
					//	spriest.Log(sim, "\thighestPossibleDmg[%d]: %01.f", i, v)
					//}
					highestPossibleIdx = i
					highestPossibleDmg = v
				}
			}
			mfAddIdx := highestPossibleIdx

			if mfAddIdx == 0 {
				chosenMfs = 0
			} else if mfAddIdx == 1 {
				chosenMfs = 1
			} else if mfAddIdx == 2 {
				chosenMfs = 2
			} else if mfAddIdx == 3 {
				chosenMfs = 3
			}
			if chosenMfs > 0 {
				if allCDs[mbIdb].Seconds() < currentWait.Seconds() && allCDs[mbIdb].Seconds() == 0 && (mfAddIdx == 2 && spellDPCT[0] > spellDPCT[4]/3*2) || (mfAddIdx == 3 && spellDPCT[0] > spellDPCT[4]) {
					bestIdx = 0
					currentWait = allCDs[mbIdb]
				} else if tickLength*3 <= gcd {
					bestIdx = 4 // TODO looks fishy, repeated bestIdx = 4
				} else {
					bestIdx = 4
				}
			}
		}

		if bestIdx == 2 && allCDs[mbIdb].Seconds() < currentWait.Seconds() && currentWait.Seconds() > 0.4 {
			bestIdx = 0
			currentWait = allCDs[mbIdb]
		}

		// if current spell is SWD and mf2 is less than GCD, and is more dps than SWD then use instead
		if bestIdx == 3 && tickLength*2 <= gcd {
			if spellDPCT[3] < spellDPCT[4]*2/3 {
				bestIdx = 4
			}
		}

		// if MF1 is chosen, and SWD is off CD and isn't 0 dmg, then use SWD unless mf2 is < gcd
		if chosenMfs == 1 && allCDs[swdIdx] == 0 && swdDamage != 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				currentWait = 0
			}
		}

		if (overwriteDPS-currDPS > 200 && (currentWait < gcd/2 || float64(currentWait) >= currDotTickSpeed*0.9)) && bestIdx != 2 {
			bestIdx = 1
			currentWait = 0
		}

		if overwriteDPS-currDPS > 200 && currentWait <= gcd && currentWait >= gcd/2 && allCDs[swdIdx] == 0 {
			if tickLength*2 <= gcd {
				bestIdx = 4
			} else {
				bestIdx = 3
				currentWait = 0
			}
		}

		// if MF2 is chosen in order to get to 5 weaving stacks, then make sure that VT/DP are already up first
		if castMf2 > 0 {
			if !spriest.DevouringPlagueDot.IsActive() && swStacks >= 4 && dpDamage != 0 {
				bestIdx = 1
			} else if !spriest.VampiricTouchDot.IsActive() && swStacks >= 4 && spriest.DevouringPlagueDot.IsActive() && vtDamage != 0 {
				bestIdx = 2
			} else {
				bestIdx = 4
			}
		}
		// if at 5 SW stacks and SWP is not up, then cast unless VT/DP are down
		if swStacks == 5 && !spriest.ShadowWordPainDot.IsActive() {
			if !spriest.DevouringPlagueDot.IsActive() && swStacks >= 4 && dpDamage != 0 {
				bestIdx = 1
			} else if !spriest.VampiricTouchDot.IsActive() && swStacks >= 4 && spriest.DevouringPlagueDot.IsActive() && vtDamage != 0 {
				bestIdx = 2
			} else {
				bestIdx = 5
			}
		}
		// cast SWP 3rd for short fights
		if castSwpNow {
			bestIdx = 5
		}
		// Snap shot BL on DP
		if overwriteDPS2-currDPS2 > 200 && bestIdx != 2 { //Seems to be a dps loss to overwrite a DP to snap shot
			bestIdx = 1
			currentWait = 0
		}

		// If BL is almost up and VT is not active, then use VT
		if deltaTimeBL <= gcd && !spriest.VampiricTouchDot.IsActive() && deltaTimeBL > 0 {
			bestIdx = 2
		}
		// If BL is up in <0.3 seconds and greater than 10ms, then wait for it to be active
		if deltaTimeBL <= time.Millisecond*300 && deltaTimeBL > time.Millisecond*10 {
			bestIdx = 1
			currentWait = time.Millisecond * time.Duration(math.Round(deltaTimeBL.Seconds()*1010))
		}

		//if sim.Log != nil {
		//spriest.Log(sim, "spriest.BLUsedAt %d", currentWait)
		//spriest.Log(sim, "dpDamage %d", dpDamage)
		//spriest.Log(sim, "currentWait %d", currentWait)
		//}
		if spriest.PrevTicks == 4 {
			castMf2 = 1
			bestIdx = 4
			spriest.PrevTicks = 0
		}

		if castMf2 == 1 && allCDs[mbIdb].Seconds() == 0 {
			bestIdx = 0
		}

		if currentWait > 0 && bestIdx != 5 && bestIdx != 4 {
			spriest.WaitUntil(sim, sim.CurrentTime+currentWait)
			return
		}

	}

	if bestIdx == 0 {
		spell = spriest.MindBlast
	} else if bestIdx == 1 {
		spell = spriest.DevouringPlague
	} else if bestIdx == 2 {
		spell = spriest.VampiricTouch
	} else if bestIdx == 3 {
		spell = spriest.ShadowWordDeath
	} else if bestIdx == 5 {
		spell = spriest.ShadowWordPain // once swp is cast need a way for talents to refresh the duration
	} else if bestIdx == 4 {

		if castMf2 == 0 {
			if spriest.InnerFocus != nil && spriest.InnerFocus.IsReady(sim) {
				spriest.InnerFocus.Cast(sim, nil)
			}
		}

		var numTicks int

		if rotType == proto.ShadowPriest_Rotation_Basic || rotType == proto.ShadowPriest_Rotation_Clipping {

			if spriest.MindBlast.TimeToReady(sim) == 0 {
				spell = spriest.MindBlast
				if success := spell.Cast(sim, spriest.CurrentTarget); !success {
					spriest.WaitForMana(sim, spell.CurCast.Cost)
				}
				return
			} else if spriest.ShadowWordDeath.TimeToReady(sim) == 0 {
				spell = spriest.ShadowWordDeath
				if success := spell.Cast(sim, spriest.CurrentTarget); !success {
					spriest.WaitForMana(sim, spell.CurCast.Cost)
				}
				return
			} else {
				if rotType == proto.ShadowPriest_Rotation_Basic {
					//numTicks = spriest.BasicMindflayRotation(sim, allCDs, gcd, tickLength)
					numTicks = 3
				} else if rotType == proto.ShadowPriest_Rotation_Clipping {
					numTicks = spriest.ClippingMindflayRotation(sim, allCDs, gcd, tickLength)
				}
			}
		} else {
			if chosenMfs == 1 {
				numTicks = 1 // determiend above that it's more dps to add MF1, need if it's not better to enter ideal rotation instead
			} else if (castMf2 == 1 && spriest.DevouringPlagueDot.IsActive() && spriest.VampiricTouchDot.IsActive()) || (deltaTimeBL < tickLength*3 && deltaTimeBL > time.Millisecond*200) {
				if spriest.MindFlayTickDuration()*3 < gcd {
					numTicks = 3
				} else {
					numTicks = 2
				}
			} else {
				numTicks = spriest.IdealMindflayRotation(sim, allCDs, gcd, tickLength, currentWait, mfDamage, mbDamage, dpDamage, vtDamage, swdDamage, overwriteDPS) //enter the mf optimizaiton routine to optimze mf clips and for next optimal spell
			}
		}

		if numTicks == 0 {
			// Means we'd rather wait for next CD (swp, vt, etc) than start a MF cast.
			nextCD := core.NeverExpires
			for _, v := range allCDs[1 : len(allCDs)-1] {
				if v < nextCD {
					nextCD = v
				}
			}
			spriest.WaitUntil(sim, sim.CurrentTime+nextCD)
			return
		}
		if numTicks == 2 && allCDs[mbIdb].Seconds() == 0 {
			spell = spriest.MindBlast
		} else {
			spell = spriest.MindFlay[numTicks]
		}
	} else {

		mbcd := spriest.MindBlast.TimeToReady(sim)
		swdcd := spriest.ShadowWordDeath.TimeToReady(sim)
		vtidx := spriest.VampiricTouchDot.RemainingDuration(sim) - vtCastTime
		swpidx := spriest.ShadowWordPainDot.RemainingDuration(sim)
		dpidx := spriest.DevouringPlagueDot.RemainingDuration(sim)
		wait1 = core.MinDuration(mbcd, swdcd)
		wait2 = core.MinDuration(dpidx, wait1)
		wait3 = core.MinDuration(vtidx, swpidx)
		wait = core.MinDuration(wait3, wait2)
		if wait <= 0 {
			spriest.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
		} else {
			spriest.WaitUntil(sim, sim.CurrentTime+wait)
		}
		return
	}
	if success := spell.Cast(sim, spriest.CurrentTarget); !success {
		spriest.WaitForMana(sim, spell.CurCast.Cost)
	}
}

// Returns the number of MF ticks to use, or 0 to wait for next CD.
func (spriest *ShadowPriest) BasicMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration) int {
	// just do MF3, never clipping
	nextCD := core.NeverExpires
	for _, v := range allCDs {
		if v < nextCD {
			nextCD = v
		}
	}
	// But don't start a MF if we can't get a single tick off.
	if nextCD < gcd {
		return 0
	} else {
		return 3
	}
}

// Returns the number of MF ticks to use, or 0 to wait for next CD.
func (spriest *ShadowPriest) IdealMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration,
	currentWait time.Duration, mfDamage, mbDamage, dpDamage, vtDamage, swdDamage, overwriteDPS float64) int {
	nextCD := core.NeverExpires
	nextIdx := -1

	newCDs := []time.Duration{
		core.MaxDuration(0, allCDs[0]),
		core.MaxDuration(0, allCDs[1]),
		core.MaxDuration(0, allCDs[2]),
	}

	for i, v := range newCDs {
		if v < nextCD {
			nextCD = v
			nextIdx = i
		}
	}

	if currentWait != 0 {
		nextCD = currentWait
	}

	var numTicks int
	numTicks_Base := 0.0
	numTicks_floored := 0.0
	if nextCD < gcd/2 {
		numTicks = 0
	} else {
		numTicks_Base = nextCD.Seconds() / tickLength.Seconds()
		numTicks_floored = math.Floor(nextCD.Seconds() / tickLength.Seconds())
		numTicks = int(numTicks_Base)
	}

	AlmostAnotherTick := numTicks_Base - numTicks_floored

	if AlmostAnotherTick > 0.95 {
		numTicks += 1
	}

	mfTickDamage := mfDamage * 0.3333

	if tickLength.Seconds() < gcd.Seconds()/2.9 {
		numTicks = 3
		return numTicks
	}

	//if sim.Log != nil {
	//spriest.Log(sim, "AlmostAnotherTick %d", AlmostAnotherTick)
	//spriest.Log(sim, "numTicks %d", numTicks)
	//spriest.Log(sim, "tickLength %d", tickLength.Seconds())
	//spriest.Log(sim, "nextCD %d", nextCD.Seconds())
	//spriest.Log(sim, "numTicks_Base %d", numTicks_Base)
	//spriest.Log(sim, "numTicks_floored %d", numTicks_floored)
	//}

	if numTicks < 100 && overwriteDPS == 0 { // if the code entered this loop because mf is the higest dps spell, and the number of ticks that can fit in the remaining cd time is < 1, then just cast a mf3 as it essentially fits perfectly
		// TODO: Should spriest latency be added to the second option here?

		mfTime := core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
		if numTicks == 0 {
			mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
		}

		if sim.Log != nil {
			//spriest.Log(sim, "mfTime %d", mfTime.Seconds())
			//spriest.Log(sim, "allCDs %d", allCDs[0].Seconds())
			//spriest.Log(sim, "mf3Time %d", float64(time.Duration(3*tickLength).Seconds()))
		}
		// Amount of gap time after casting mind flay, but before each CD is available.

		cdDiffs := []time.Duration{
			core.MaxDuration(0, allCDs[0]-mfTime),
			core.MaxDuration(0, allCDs[1]-mfTime),
			core.MaxDuration(0, allCDs[2]-mfTime),
			core.MaxDuration(0, allCDs[3]-mfTime),
			0,
		}

		mfspdmg := 0.0
		if numTicks != 0 {
			mfspdmg = mfTickDamage * float64(numTicks) / (time.Duration(numTicks) * tickLength).Seconds()
		} else if numTicks > 3 {
			mfspdmg = mfTickDamage * float64(3) / (time.Duration(3) * tickLength).Seconds()
		}
		if sim.Log != nil {
			//spriest.Log(sim, "mfspdmg %d", mfspdmg)
		}
		spellDamages := []float64{
			// MB dps
			mbDamage / (gcd + cdDiffs[mbIdb]).Seconds(),
			// DP dps
			dpDamage / (gcd + cdDiffs[dpIdx]).Seconds(),
			// VT dps
			vtDamage / (gcd + cdDiffs[vtIdx]).Seconds(),
			// SWD dps
			swdDamage / (gcd + cdDiffs[swdIdx]).Seconds(),

			mfspdmg,
		}

		bestIdx := 0
		bestDmg := 0.0
		for i, v := range spellDamages {
			if sim.Log != nil {
				//spriest.Log(sim, "\tspellDamages[%d]: %01.f", i, v)
			}
			if v > bestDmg {
				bestIdx = i
				bestDmg = v
			}
		}

		//if numTicks < 1 && bestIdx == 4 {
		//	numTicks = 3
		//return numTicks
		//}

		if sim.Log != nil {
			spriest.Log(sim, "bestIdx %d", bestIdx)
			spriest.Log(sim, "nextIdx %d", nextIdx)
			spriest.Log(sim, "spellDamages[bestIdx]  %d", spellDamages[bestIdx])
			spriest.Log(sim, "spellDamages[nextIdx]  %d", spellDamages[nextIdx])
			spriest.Log(sim, "numTicks %d", numTicks)
		}

		if bestIdx != nextIdx && spellDamages[nextIdx] < spellDamages[bestIdx] && bestIdx != 4 {
			numTicks_Base = allCDs[bestIdx].Seconds() / tickLength.Seconds()
			numTicks_floored = math.Floor(allCDs[bestIdx].Seconds() / tickLength.Seconds())
			numTicks = int(numTicks_Base)
			if sim.Log != nil {
				spriest.Log(sim, "numTicks2 %d", numTicks)
			}
			AlmostAnotherTick := numTicks_Base - numTicks_floored

			if AlmostAnotherTick > 0.75 {
				numTicks += 1
			}

			mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength)
			if numTicks > 3 && numTicks < 5 {
				addedgcd := core.MaxDuration(gcd, time.Duration(2)*tickLength)
				addedgcdtime := addedgcd - time.Duration(2)*tickLength
				mfTime = core.MaxDuration(gcd, time.Duration(numTicks)*tickLength+2*addedgcdtime)
			}
			deltaTime := allCDs[bestIdx] - mfTime
			cdDiffs = []time.Duration{
				core.MaxDuration(0, allCDs[0]-mfTime),
				core.MaxDuration(0, allCDs[1]-mfTime),
				core.MaxDuration(0, allCDs[2]-mfTime),
				core.MaxDuration(0, allCDs[3]-mfTime),
				0,
			}
			if deltaTime.Seconds() < -0.33 {
				numTicks -= 1
				cdDiffs[bestIdx] += tickLength
			}
		}

		if numTicks < 0 {
			numTicks = 0
		}

		chosenWait := cdDiffs[bestIdx]

		//if sim.Log != nil {
		//spriest.Log(sim, "numTicks %d", numTicks)
		//spriest.Log(sim, "mfTime %d", mfTime.Seconds())
		//spriest.Log(sim, "chosenWait %d", chosenWait.Seconds())
		//}

		var newInd int
		if chosenWait > gcd {
			check_CDs := cdDiffs
			check_CDs[bestIdx] = time.Second * 15
			// need to find a way to sort the cdDiffs and find the next highest dps cast with lower wait time
			for i, v := range check_CDs {
				if v < nextCD {
					//nextCDc = v
					newInd = i
				}
			}
		}
		skipNext := 0
		var totalWaitCurr float64
		var numTicksAvail float64
		var remainTime1 float64
		var remainTime2 float64
		var remainTime3 float64
		var addTime1 float64
		var addTime2 float64
		var addTime3 float64

		if chosenWait.Seconds() > gcd.Seconds() && bestIdx != newInd && newInd > -1 {

			tick_var := float64(numTicks)
			if numTicks == 1 {
				totalWaitCurr = chosenWait.Seconds() - tick_var*gcd.Seconds()
			} else {
				totalWaitCurr = chosenWait.Seconds() - tick_var*tickLength.Seconds()
			}

			if totalWaitCurr-gcd.Seconds() <= gcd.Seconds() {
				if totalWaitCurr > tickLength.Seconds() {
					numTicksAvail = math.Floor((totalWaitCurr - gcd.Seconds()) / tickLength.Seconds())
				} else {
					numTicksAvail = math.Floor((totalWaitCurr - gcd.Seconds()) / gcd.Seconds())
				}
			} else {
				numTicksAvail = math.Floor((totalWaitCurr - gcd.Seconds()) / tickLength.Seconds())
			}

			if numTicksAvail < 0 {
				numTicksAvail = 0
			}

			// TODO looks fishy, remainTime1 and remainTime2 are equal
			remainTime1 = totalWaitCurr - tickLength.Seconds()*numTicksAvail - gcd.Seconds()
			remainTime2 = totalWaitCurr - 1*tickLength.Seconds()*numTicksAvail - gcd.Seconds()
			remainTime3 = totalWaitCurr - 2*tickLength.Seconds()*numTicksAvail - gcd.Seconds()

			if remainTime1 > 0 {
				addTime1 = remainTime1
			} else {
				addTime1 = 0
			}

			if remainTime2 > 0 {
				addTime2 = remainTime2
			} else {
				addTime2 = 0
			}

			if remainTime3 > 0 {
				addTime3 = remainTime3
			} else {
				addTime3 = 0
			}

			dpsPossible0 := []float64{
				0,
				0,
				0,
			}

			cd_dpsb := spellDamages[bestIdx]
			cd_dpsn := spellDamages[newInd]

			dpsPossible0[0] = (numTicksAvail*mfTickDamage + cd_dpsb*gcd.Seconds() + cd_dpsn*gcd.Seconds()) / (numTicksAvail*tickLength.Seconds() + 2*gcd.Seconds() + addTime1)
			dpsPossible0[1] = (tick_var*mfTickDamage + cd_dpsb*(cdDiffs[bestIdx].Seconds()+gcd.Seconds()) + cd_dpsn*(cdDiffs[newInd].Seconds())) / (tick_var*tickLength.Seconds() + (cdDiffs[bestIdx].Seconds() + gcd.Seconds()) + (cdDiffs[newInd].Seconds() + addTime2))
			dpsPossible0[2] = ((tick_var+1)*mfTickDamage + cd_dpsb*(cdDiffs[len(cdDiffs)-1-1].Seconds()+gcd.Seconds()) + cd_dpsn*(cdDiffs[len(cdDiffs)-1].Seconds()-tickLength.Seconds())) / ((tick_var+1)*tickLength.Seconds() + (cdDiffs[bestIdx].Seconds() + gcd.Seconds()) + (cdDiffs[newInd].Seconds() + addTime3))

			highestPossibleDmg := 0.0
			highestPossibleIdx := -1
			// TODO looks fishy, this branch is never taken
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible0 {

					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}
			if highestPossibleIdx > 0 {
				numTicks = highestPossibleIdx + 1
			} else {
				numTicks = int(numTicksAvail)
				skipNext = 1
			}
		}

		if sim.Log != nil {
			spriest.Log(sim, "numTicks3 %d", numTicks)
		}

		if numTicks > 3 {
			if (allCDs[bestIdx] - time.Duration(numTicks-1)*tickLength - gcd) >= 0 {
				//if (allCDs[3]-time.Duration(numTicks-1)*tickLength <= 0) || (allCDs[0]-time.Duration(numTicks-1)*tickLength <= 0) { \\might need to readd this for later phases
				if allCDs[3]-time.Duration(numTicks-1)*tickLength <= 0 {
					numTicks = 3
					return numTicks
				}
			}
		}

		if skipNext == 0 {
			finalMFStart := math.Mod(float64(numTicks), 3) // Base ticks before adding additional
			dpsPossible := []float64{
				bestDmg, // dps with no tick and just wait
				0,
				0,
				0,
			}
			dpsDuration := (chosenWait + gcd).Seconds()

			highestPossibleIdx := 0
			// TODO: Modified this slightly to expand time window, but it still doesn't change dps for any tests.
			// Probably can remove this entirely (and then also the if highestPossibleIdx == 0 right after)
			if (finalMFStart == 2) && (chosenWait <= tickLength && chosenWait > (tickLength-time.Millisecond*15)) {
				highestPossibleIdx = 1 // if the wait time is equal to an extra mf tick, and there are already 2 ticks, then just add 1
			}

			if highestPossibleIdx == 0 {
				switch finalMFStart {
				case 0:
					// this means that the extra ticks will be relative to starting a new mf cast entirely
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / (gcd.Seconds() + gcd.Seconds())          // new damage for 1 extra tick
					dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / (2*tickLength.Seconds() + gcd.Seconds()) // new damage for 2 extra tick
					dpsPossible[3] = (bestDmg*dpsDuration + mfDamage) / (3*tickLength.Seconds() + gcd.Seconds())     // new damage for 3 extra tick

				case 1:
					total_check_time := 2 * tickLength

					if total_check_time < gcd {
						newDuration := (gcd + gcd).Seconds()
						dpsPossible[1] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * (finalMFStart + 1))) / newDuration
					} else {
						newDuration := ((total_check_time - gcd) + gcd).Seconds()
						dpsPossible[1] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * (finalMFStart + 1))) / newDuration
					}
					// % check add 2
					total_check_time2 := 2 * tickLength.Seconds()
					if total_check_time2 < gcd.Seconds() {
						dpsPossible[2] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * (finalMFStart + 2))) / (gcd.Seconds() + gcd.Seconds())
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + (mfDamage * 1 / 3 * (finalMFStart + 2))) / (total_check_time2 + gcd.Seconds())
					}
				case 2:
					// % check add 1
					total_check_time := tickLength
					newDuration := (total_check_time + gcd).Seconds()
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / newDuration

				default:
					dpsPossible[1] = (bestDmg*dpsDuration + mfDamage*1/3) / (gcd.Seconds() + gcd.Seconds())
					if tickLength*2 > gcd {
						dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / (2*tickLength.Seconds() + gcd.Seconds())
					} else {
						dpsPossible[2] = (bestDmg*dpsDuration + mfDamage*2/3) / (gcd.Seconds() + gcd.Seconds())
					}
					dpsPossible[3] = (bestDmg*dpsDuration + mfDamage) / (3*tickLength.Seconds() + gcd.Seconds())
				}
			}

			// Find the highest possible dps and its index
			// highestPossibleIdx := 0
			highestPossibleDmg := 0.0
			if highestPossibleIdx == 0 {
				for i, v := range dpsPossible {
					if sim.Log != nil {
						//spriest.Log(sim, "\tdpsPossible[%d]: %01.f", i, v)
					}
					if v >= highestPossibleDmg {
						highestPossibleIdx = i
						highestPossibleDmg = v
					}
				}
			}

			numTicks += highestPossibleIdx
			if sim.Log != nil {
				spriest.Log(sim, "numTicks4 %d", numTicks)
			}
			// if sim.Log != nil {
			// 	spriest.Log(sim, "final_ticks %d", numTicks)
			// }
			if numTicks == 1 && tickLength*3 <= time.Duration(float64(gcd)*1.05) {
				numTicks += 2
			}
			if numTicks == 1 && tickLength*2 <= gcd {
				numTicks += 1
			}
			//  Now that the number of optimal ticks has been determined to optimize dps
			//  Now optimize mf2s and mf3s

			//if numTicks == 0 {
			//return numTicks
			//}

			if numTicks == 1 {
				numTicks = 1
			} else if numTicks == 2 || numTicks == 4 {
				if numTicks == 4 {
					spriest.PrevTicks = 4
				}
				numTicks = 2
			} else if numTicks == 0 {
				numTicks = 2
			} else {
				numTicks = 3
			}
		}
	} else {
		numTicks = int(nextCD / tickLength)
		if nextCD-core.MaxDuration(gcd, time.Duration(2)*tickLength) < 0 && numTicks != 0 {
			numTicks -= 1
		}
		// if sim.Log != nil {
		// 	spriest.Log(sim, "c_ticks %d", numTicks)
		// 	spriest.Log(sim, "nextCD %d", nextCD)
		// 	spriest.Log(sim, "tickLength %d", tickLength)
		// }
		if numTicks == 0 {
			// if sim.Log != nil {
			//   spriest.Log(sim, "zero ticks %d", numTicks)
			// }
			numTicks = 2
		}
		if numTicks >= 3 {
			numTicks = 3
		}
	}

	return numTicks
}

func (spriest *ShadowPriest) ClippingMindflayRotation(sim *core.Simulation, allCDs []time.Duration, gcd time.Duration, tickLength time.Duration) int {
	nextCD := core.NeverExpires
	for _, v := range allCDs[1 : len(allCDs)-1] {
		if v < nextCD {
			nextCD = v
		}
	}

	// if sim.Log != nil {
	// 	spriest.Log(sim, "<spriest> NextCD: %0.2f", nextCD.Seconds())
	// }

	// This means a CD is coming up before we could cast a single MF
	if nextCD < gcd {
		return 0
	}

	// How many ticks we have time for.
	numTicks := int((nextCD - time.Duration(spriest.rotation.Latency)) / tickLength)
	if numTicks == 1 {
		return 1
	} else if numTicks == 2 || numTicks == 4 {
		return 2
	} else {
		return 3
	}
}
