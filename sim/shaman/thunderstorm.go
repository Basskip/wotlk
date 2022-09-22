package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var thunderstormActionID = core.ActionID{SpellID: 59159}

// newThunderstormSpell returns a precomputed instance of lightning bolt to use for casting.
func (shaman *Shaman) registerThunderstormSpell() {
	if !shaman.Talents.Thunderstorm {
		return
	}

	manaRestore := 0.08
	if shaman.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfThunderstorm) {
		manaRestore = 0.1
	}

	manaMetrics := shaman.NewManaMetrics(thunderstormActionID)

	spellConfig := core.SpellConfig{
		ActionID:    thunderstormActionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		BonusCritRating:  core.TernaryFloat64(shaman.Talents.CallOfThunder, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)),
		ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),

		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics, true)
		},
	}

	if shaman.thunderstormInRange {
		effect := core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMagic(1450, 1656, 0.172),
			OutcomeApplier: shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
		}
		aoeApply := core.ApplyEffectFuncAOEDamageCapped(shaman.Env, effect)
		spellConfig.ApplyEffects = func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			aoeApply(sim, unit, spell)                                           // Calculates hits/crits/dmg on each target
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics, true) // adds mana no matter what
		}
	}
	shaman.Thunderstorm = shaman.RegisterSpell(spellConfig)
}
