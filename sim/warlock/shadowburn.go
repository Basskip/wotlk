package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerShadowBurnSpell() {
	baseCost := 0.2 * warlock.BaseMana
	actionID := core.ActionID{SpellID: 47827}
	spellSchool := core.SpellSchoolShadow

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowburn) {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
				if isExecute == 35 {
					warlock.Shadowburn.BonusCritRating += 20 * core.CritRatingPerCritChance
				}
			})
		})
	}

	effect := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMagic(775.0, 865.0, 0.429*(1+0.04*float64(warlock.Talents.ShadowAndFlame))),
		OutcomeApplier: warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
	}

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:  core.GCDDefault, // backdraft procs don't change the GCD of shadowburn
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}
