package component

import (
	"mob/pkg/timer"

	"github.com/sedyh/mizu/pkg/engine"
)

type ModType int

const (
	Attack ModType = iota
	Buff
	Debuff
)

type ModTarget int

const (
	TargetAlly ModTarget = iota
	TargetEnemy
	TargetSelf
)

type ModOrder int

const (
	Normal ModOrder = iota
	BeforeAll
	AfterAll
	BeforeAttack
	AfterAttack
	OrderRandom
)

type Mod struct {
	Name   string
	Desc   string
	Type   ModType
	Target ModTarget
	Range  NPCRange
	Order  ModOrder
	Move   Move
	Once   bool
}

func (m *Mod) IsGood() bool {
	switch m.Type {
	case Attack, Debuff:
		return m.Target == TargetEnemy
	case Buff:
		return m.Target == TargetAlly || m.Target == TargetSelf
	default:
		return false
	}
}

type NPCType int

const (
	Ally NPCType = iota
	Enemy
)

type NPCStrategy int

const (
	Closest NPCStrategy = iota
	Farthest
	HighAttack
	LowAttack
	HighHealth
	LowHealth
	Random
)

type NPCRange float64

const (
	Melee  NPCRange = 30
	Ranged NPCRange = 90
)

type NPCAnimation float64

const (
	Neutral NPCAnimation = iota
	Happy
	Sad
	Angry
	Spinning
	Bouncing
	Walking
	Falling
	EyesClosed
)

type NPCAnimationSprite struct {
	Part      NPCAnimationPart
	Animation NPCAnimation
	Path      string
	Frames    int
}

func NewNPCAnimationSprite(opts ...NPCAnimationSpriteOption) (a NPCAnimationSprite) {
	for _, opt := range opts {
		opt(&a)
	}
	return
}

type NPCAnimationSpriteOption func(*NPCAnimationSprite)

func WithPart(part NPCAnimationPart) NPCAnimationSpriteOption {
	return func(ns *NPCAnimationSprite) {
		ns.Part = part
	}
}

func WithAnimation(animation NPCAnimation) NPCAnimationSpriteOption {
	return func(ns *NPCAnimationSprite) {
		ns.Animation = animation
	}
}

func WithPath(path string) NPCAnimationSpriteOption {
	return func(ns *NPCAnimationSprite) {
		ns.Path = path
	}
}

type NPCAnimationPart int

const (
	AnimFill NPCAnimationPart = iota
	AnimBase
	AnimEyes
	AnimArm
	AnimFeet
)

type NPC struct {
	Mods []Mod
	// mods gained during combat
	combatMods []Mod
	Name       string
	Type       NPCType
	// targeting strategy
	Strategy NPCStrategy
	// center animation sprites?
	Center           bool
	Animations       map[NPCAnimation]bool
	AnimationSprites map[NPCAnimationPart]NPCAnimationSprite
	AnimationTimer   timer.Timer
}

func NewNPC(options ...NPCOption) (n NPC) {
	n.Animations = make(map[NPCAnimation]bool)
	n.AnimationSprites = make(map[NPCAnimationPart]NPCAnimationSprite)
	n.AnimationTimer = timer.NewTimer(1)
	n.AddAnimations(Neutral)
	for _, opt := range options {
		opt(&n)
	}
	return
}

type NPCOption func(*NPC)

func WithMod(m Mod) NPCOption {
	return func(n *NPC) {
		n.AddMod(m)
	}
}

func WithMods(mods []Mod) NPCOption {
	return func(n *NPC) {
		for _, m := range mods {
			n.AddMod(m)
		}
	}
}

func WithName(name string) NPCOption {
	return func(n *NPC) {
		n.Name = name
	}
}

func WithType(t NPCType) NPCOption {
	return func(n *NPC) {
		n.Type = t
	}
}

func (n *NPC) AddMod(m Mod) {
	var mods []Mod
	found := false
	for _, am := range n.Mods {
		if am.Name == m.Name {
			found = true
		} else {
			mods = append(mods, am)
		}
	}
	if !found {
		mods = append(mods, m)
	}
	n.Mods = mods
}

func (n *NPC) AddCombatMod(m Mod) {
	n.combatMods = append(n.combatMods, m)
}

func (n *NPC) ClearCombatMods(m Mod) {
	n.combatMods = make([]Mod, 0)
}

func (n *NPC) RemoveMod(m Mod) {
	var mods []Mod
	for _, am := range n.Mods {
		if am.Name != m.Name {
			mods = append(mods, m)
		}
	}
	n.Mods = mods
}

func (n *NPC) AddAnimations(animations ...NPCAnimation) {
	if n.Animations == nil {
		n.Animations = make(map[NPCAnimation]bool)
	}
	for _, anim := range animations {
		n.Animations[anim] = true
	}
}

func (n *NPC) RemoveAnimations(animations ...NPCAnimation) {
	if n.Animations == nil {
		n.Animations = make(map[NPCAnimation]bool)
	}
	for _, anim := range animations {
		n.Animations[anim] = false
	}
}

func (n *NPC) UseAnimationSprite(animation NPCAnimationSprite) {
	if n.AnimationSprites == nil {
		n.AnimationSprites = make(map[NPCAnimationPart]NPCAnimationSprite)
	}
	existing, ok := n.AnimationSprites[animation.Part]
	if !ok || existing.Path != animation.Path {
		n.AnimationSprites[animation.Part] = animation
	}
}

// update sprites for current animations
func (n *NPC) UpdateAnimations() {
	atLeastOne := false
	// for anim := range n.Animations {
	// 	atLeastOne = true
	// 	// switch anim {
	// 	// default:

	// 	// }
	// }
	if !atLeastOne {
		// neutral
		n.UseAnimationSprite(NewNPCAnimationSprite(
			WithAnimation(Neutral),
			WithPart(AnimFill),
			WithPath("asset/image/ally/fill.png"),
		))
		n.UseAnimationSprite(NewNPCAnimationSprite(
			WithAnimation(Neutral),
			WithPart(AnimBase),
			WithPath("asset/image/ally/base.png"),
		))
		n.UseAnimationSprite(NewNPCAnimationSprite(
			WithAnimation(Neutral),
			WithPart(AnimEyes),
			WithPath("asset/image/ally/eyes.png"),
		))
		n.UseAnimationSprite(NewNPCAnimationSprite(
			WithAnimation(Neutral),
			WithPart(AnimArm),
			WithPath("asset/image/ally/arm.png"),
		))
		n.UseAnimationSprite(NewNPCAnimationSprite(
			WithAnimation(Neutral),
			WithPart(AnimFeet),
			WithPath("asset/image/ally/feet.png"),
		))
	}
}

func GetAllNPCOfType(w engine.World, npcType NPCType) (entities []engine.Entity) {
	for _, e := range w.View(NPC{}).Filter() {
		var npc *NPC
		e.Get(&npc)
		if npc.Type == npcType {
			entities = append(entities, e)
		}
	}
	return
}
