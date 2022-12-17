import{A as t,cT as e,cU as n,cV as a,cW as i,bO as s,P as o,j as l,cX as r,aZ as d,bv as m,cY as c,G as u,bA as f,bD as S,bF as g,bp as p,d as h,a2 as b,a as v,T as O,ai as y,aJ as P,ag as I,F as T,aj as A}from"./raid_sim_action-2e58abdd.chunk.js";import{k,b as H,c as W,z as D,i as w,j as R,B as E,I as B,a1 as M,T as x,U as C,V as j,W as N,a2 as G,s as J}from"./individual_sim_ui-35391952.chunk.js";const F=k({fieldName:"customRotation",numColumns:2,values:[{actionId:t.fromSpellId(53408),value:e.JudgementOfWisdom},{actionId:t.fromSpellId(48806),value:e.HammerOfWrath},{actionId:t.fromSpellId(48819),value:e.Consecration},{actionId:t.fromSpellId(48817),value:e.HolyWrath},{actionId:t.fromSpellId(48801),value:e.Exorcism},{actionId:t.fromSpellId(61411),value:e.ShieldOfRighteousness},{actionId:t.fromSpellId(48827),value:e.AvengersShield},{actionId:t.fromSpellId(53595),value:e.HammerOfTheRighteous},{actionId:t.fromSpellId(48952),value:e.HolyShield}]}),V={inputs:[H({fieldName:"hammerFirst",label:"Open with HotR",labelTooltip:"Open with Hammer of the Righteous instead of Shield of Righteousness in the standard rotation. Recommended for AoE."}),H({fieldName:"squeezeHolyWrath",label:"Squeeze Holy Wrath",labelTooltip:"Squeeze a Holy Wrath cast during sufficiently hasted GCDs (Bloodlust) in the standard rotation."}),W({fieldName:"waitSlack",label:"Max Wait Time (ms)",labelTooltip:"Maximum time in milliseconds to prioritize waiting for next Hammer/Shield to maintain 969. Affects standard and custom priority."}),H({fieldName:"useCustomPrio",label:"Use custom priority",labelTooltip:"Deviates from the standard 96969 rotation, using the priority configured below. Will still attempt to keep a filler GCD between Hammer and Shield."}),F]},z=D({fieldName:"aura",label:"Aura",values:[{name:"None",value:n.NoPaladinAura},{name:"Devotion Aura",value:n.DevotionAura},{name:"Retribution Aura",value:n.RetributionAura}]}),U=D({fieldName:"seal",label:"Seal",labelTooltip:"The seal active before encounter",values:[{name:"Vengeance",value:a.Vengeance},{name:"Command",value:a.Command}]}),L=D({fieldName:"judgement",label:"Judgement",labelTooltip:"Judgement debuff you will use on the target during the encounter.",values:[{name:"Wisdom",value:i.JudgementOfWisdom},{name:"Light",value:i.JudgementOfLight}]}),q=w({fieldName:"useAvengingWrath",label:"Use Avenging Wrath"}),_=R({fieldName:"damageTakenPerSecond",label:"Damage Taken Per Second",labelTooltip:"Damage taken per second across the encounter. Used to model mana regeneration from Spiritual Attunement. This value should NOT include damage taken from Seal of Blood / Judgement of Blood. Leave at 0 if you do not take damage during the encounter."}),K={name:"Baseline Example",data:s.create({talentsString:"-05005135200132311333312321-511302012003",glyphs:{major1:o.GlyphOfSealOfVengeance,major2:o.GlyphOfRighteousDefense,major3:o.GlyphOfDivinePlea,minor1:l.GlyphOfSenseUndead,minor2:l.GlyphOfLayOnHands,minor3:l.GlyphOfBlessingOfKings}})},X=r.create({hammerFirst:!1,squeezeHolyWrath:!0,waitSlack:300,useCustomPrio:!1,customRotation:d.create({spells:[m.create({spell:e.ShieldOfRighteousness}),m.create({spell:e.HammerOfTheRighteous}),m.create({spell:e.HolyShield}),m.create({spell:e.HammerOfWrath}),m.create({spell:e.Consecration}),m.create({spell:e.AvengersShield}),m.create({spell:e.JudgementOfWisdom}),m.create({spell:e.Exorcism})]})}),Y=c.create({aura:n.RetributionAura,judgement:i.JudgementOfWisdom,damageTakenPerSecond:0}),Z=u.create({flask:f.FlaskOfStoneblood,food:S.FoodDragonfinFilet,defaultPotion:g.IndestructiblePotion,prepopPotion:g.IndestructiblePotion}),Q={name:"Preraid Preset",tooltip:E,enableWhen:t=>!0,gear:p.fromJsonString('{"items": [\n\t\t{\n\t\t\t"id": 42549,\n\t\t\t"enchant": 3818,\n\t\t\t"gems": [\n\t\t\t\t41396,\n\t\t\t\t49110\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40679\n\t\t},\n\t\t{\n\t\t\t"id": 37635,\n\t\t\t"enchant": 3852,\n\t\t\t"gems": [\n\t\t\t\t40015\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 44188,\n\t\t\t"enchant": 3605\n\t\t},\n\t\t{\n\t\t\t"id": 39638,\n\t\t\t"enchant": 1953,\n\t\t\t"gems": [\n\t\t\t\t36767,\n\t\t\t\t40089\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 37682,\n\t\t\t"enchant": 3850,\n\t\t\t"gems": [\n\t\t\t\t0\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 39639,\n\t\t\t"enchant": 3860,\n\t\t\t"gems": [\n\t\t\t\t36767,\n\t\t\t\t0\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 37379,\n\t\t\t"enchant": 3601,\n\t\t\t"gems": [\n\t\t\t\t40022,\n\t\t\t\t40008\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 37292,\n\t\t\t"enchant": 3822,\n\t\t\t"gems": [\n\t\t\t\t40089\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 44243,\n\t\t\t"enchant": 3606\n\t\t},\n\t\t{\n\t\t\t"id": 37186\n\t\t},\n\t\t{\n\t\t\t"id": 37257\n\t\t},\n\t\t{\n\t\t\t"id": 44063,\n\t\t\t"gems": [\n\t\t\t\t36767,\n\t\t\t\t40015\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 37220\n\t\t},\n\t\t{\n\t\t\t"id": 37179,\n\t\t\t"enchant": 2673\n\t\t},\n\t\t{\n\t\t\t"id": 43085,\n\t\t\t"enchant": 3849\n\t\t},\n\t\t{\n\t\t\t"id": 40707\n\t\t}\n\t]}')},$={name:"P1 Preset",tooltip:E,enableWhen:t=>!0,gear:p.fromJsonString('{"items": [\n\t\t{\n\t\t\t"id": 40581,\n\t\t\t"enchant": 3818,\n\t\t\t"gems": [\n\t\t\t\t41380,\n\t\t\t\t36767\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40387\n\t\t},\n\t\t{\n\t\t\t"id": 40584,\n\t\t\t"enchant": 3852,\n\t\t\t"gems": [\n\t\t\t\t40008\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40410,\n\t\t\t"enchant": 3605\n\t\t},\n\t\t{\n\t\t\t"id": 40579,\n\t\t\t"enchant": 3832,\n\t\t\t"gems": [\n\t\t\t\t36767,\n\t\t\t\t40022\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 39764,\n\t\t\t"enchant": 3850,\n\t\t\t"gems": [\n\t\t\t\t0\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40580,\n\t\t\t"enchant": 3860,\n\t\t\t"gems": [\n\t\t\t\t40008,\n\t\t\t\t0\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 39759,\n\t\t\t"enchant": 3601,\n\t\t\t"gems": [\n\t\t\t\t40008,\n\t\t\t\t40008\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40589,\n\t\t\t"enchant": 3822\n\t\t},\n\t\t{\n\t\t\t"id": 39717,\n\t\t\t"enchant": 3606,\n\t\t\t"gems": [\n\t\t\t\t40089\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 40718\n\t\t},\n\t\t{\n\t\t\t"id": 40107\n\t\t},\n\t\t{\n\t\t\t"id": 44063,\n\t\t\t"gems": [\n\t\t\t\t36767,\n\t\t\t\t40089\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t"id": 37220\n\t\t},\n\t\t{\n\t\t\t"id": 40345,\n\t\t\t"enchant": 3788\n\t\t},\n\t\t{\n\t\t\t"id": 40400,\n\t\t\t"enchant": 3849\n\t\t},\n\t\t{\n\t\t\t"id": 40707\n\t\t}\n\t]}')};class tt extends B{constructor(t,e){super(t,e,{cssClass:"protection-paladin-sim-ui",knownIssues:[],epStats:[h.StatStamina,h.StatStrength,h.StatAgility,h.StatAttackPower,h.StatMeleeHit,h.StatSpellHit,h.StatMeleeCrit,h.StatExpertise,h.StatMeleeHaste,h.StatArmorPenetration,h.StatSpellPower,h.StatArmor,h.StatDefense,h.StatBlock,h.StatBlockValue,h.StatDodge,h.StatParry,h.StatResilience],epPseudoStats:[b.PseudoStatMainHandDps],epReferenceStat:h.StatSpellPower,displayStats:[h.StatHealth,h.StatArmor,h.StatStamina,h.StatStrength,h.StatAgility,h.StatAttackPower,h.StatMeleeHit,h.StatMeleeCrit,h.StatMeleeHaste,h.StatExpertise,h.StatArmorPenetration,h.StatSpellPower,h.StatSpellHit,h.StatDefense,h.StatBlock,h.StatBlockValue,h.StatDodge,h.StatParry,h.StatResilience],modifyDisplayStats:t=>{let e=new v;return O.freezeAllAndDo((()=>{t.getMajorGlyphs().includes(o.GlyphOfSealOfVengeance)&&t.getSpecOptions().seal==a.Vengeance&&(e=e.addStat(h.StatExpertise,10*M))})),{talents:e}},defaults:{gear:$.gear,epWeights:v.fromMap({[h.StatArmor]:.07,[h.StatStamina]:1.14,[h.StatStrength]:1,[h.StatAgility]:.62,[h.StatAttackPower]:.26,[h.StatExpertise]:.69,[h.StatMeleeHit]:.79,[h.StatMeleeCrit]:.3,[h.StatMeleeHaste]:.17,[h.StatArmorPenetration]:.04,[h.StatSpellPower]:.13,[h.StatBlock]:.52,[h.StatBlockValue]:.28,[h.StatDodge]:.46,[h.StatParry]:.61,[h.StatDefense]:.54},{[b.PseudoStatMainHandDps]:3.33}),consumes:Z,rotation:X,talents:K.data,specOptions:Y,raidBuffs:y.create({giftOfTheWild:P.TristateEffectImproved,powerWordFortitude:P.TristateEffectImproved,strengthOfEarthTotem:P.TristateEffectImproved,arcaneBrilliance:!0,unleashedRage:!0,leaderOfThePack:P.TristateEffectRegular,icyTalons:!0,totemOfWrath:!0,demonicPact:500,swiftRetribution:!0,moonkinAura:P.TristateEffectRegular,sanctifiedRetribution:!0,manaSpringTotem:P.TristateEffectRegular,bloodlust:!0,thorns:P.TristateEffectImproved,devotionAura:P.TristateEffectImproved,shadowProtection:!0}),partyBuffs:I.create({}),individualBuffs:T.create({blessingOfKings:!0,blessingOfSanctuary:!0,blessingOfWisdom:P.TristateEffectImproved,blessingOfMight:P.TristateEffectImproved}),debuffs:A.create({judgementOfWisdom:!0,judgementOfLight:!0,misery:!0,faerieFire:P.TristateEffectImproved,ebonPlaguebringer:!0,totemOfWrath:!0,shadowMastery:!0,bloodFrenzy:!0,mangle:!0,exposeArmor:!0,sunderArmor:!0,vindication:!0,thunderClap:P.TristateEffectImproved,insectSwarm:!0})},playerIconInputs:[],rotationInputs:V,includeBuffDebuffInputs:[],excludeBuffDebuffInputs:[],otherInputs:{inputs:[x,C,j,N,G,z,q,L,U,_,J]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[K],gear:[Q,$]}})}}export{X as D,K as G,tt as P,Y as a,Z as b,$ as c};
//# sourceMappingURL=sim-09b59b34.chunk.js.map
