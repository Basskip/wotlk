package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerCircleOfHealingSpell() {
	if !priest.Talents.CircleOfHealing {
		return
	}

	baseCost := .21 * priest.BaseMana

	baseEffect := core.SpellEffect{
		IsHealing:      true,
		BaseDamage:     core.BaseDamageConfigHealing(958, 1058, 0.4029),
		OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultHealingCritMultiplier()),
	}

	var effects []core.SpellEffect
	targets := priest.Env.Raid.GetFirstNPlayersOrPets(5 + core.TernaryInt32(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfCircleOfHealing), 1, 0))
	for _, target := range targets {
		effect := baseEffect
		effect.Target = target
		effects = append(effects, effect)
	}

	priest.CircleOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48089},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellHealing,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolytesRaiment, 4), 1.1, 1),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
