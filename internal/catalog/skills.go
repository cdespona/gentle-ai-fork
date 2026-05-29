package catalog

import "github.com/gentleman-programming/gentle-ai/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// SDD skills
	{ID: model.SkillSDDInit, Name: "sdd-init", Category: "sdd", Priority: "p0"},

	{ID: model.SkillSDDApply, Name: "sdd-apply", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDVerify, Name: "sdd-verify", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDExplore, Name: "sdd-explore", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDPropose, Name: "sdd-propose", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDSpec, Name: "sdd-spec", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDDesign, Name: "sdd-design", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDTasks, Name: "sdd-tasks", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDArchive, Name: "sdd-archive", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDOnboard, Name: "sdd-onboard", Category: "sdd", Priority: "p0"},
	// Foundation skills
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillTDD, Name: "tdd", Category: "testing", Priority: "p0"},
	{ID: model.SkillTestTypeClass, Name: "test-type-classification", Category: "testing", Priority: "p0"},
	{ID: model.SkillCaveman, Name: "caveman", Category: "productivity", Priority: "p0"},
	{ID: model.SkillCodeComments, Name: "code-comments", Category: "engineering", Priority: "p0"},
	{ID: model.SkillHexagonalArch, Name: "hexagonal-architecture", Category: "architecture", Priority: "p0"},
	{ID: model.SkillJavaDevelopment, Name: "java-development", Category: "language", Priority: "p0"},
	{ID: model.SkillKotlinDevelopment, Name: "kotlin-development", Category: "language", Priority: "p0"},
	{ID: model.SkillTacticalDDD, Name: "tactical-ddd", Category: "domain-modeling", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillJudgmentDay, Name: "judgment-day", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
	{ID: model.SkillSkillRegistry, Name: "skill-registry", Category: "workflow", Priority: "p0"},
	// Sustainable review skills
	{ID: model.SkillChainedPR, Name: "chained-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCognitiveDoc, Name: "cognitive-doc-design", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCommentWriter, Name: "comment-writer", Category: "workflow", Priority: "p0"},
	{ID: model.SkillWorkUnitCommits, Name: "work-unit-commits", Category: "workflow", Priority: "p0"},
	{ID: model.SkillMemoryRecall, Name: "memory-recall", Category: "memory", Priority: "p0"},
	{ID: model.SkillMemoryCapture, Name: "memory-capture", Category: "memory", Priority: "p0"},
	{ID: model.SkillMemoryConsolidate, Name: "memory-consolidate", Category: "memory", Priority: "p0"},
	{ID: model.SkillMemoryHandoff, Name: "memory-handoff", Category: "memory", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
