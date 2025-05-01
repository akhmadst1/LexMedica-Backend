package pkg

import (
	"fmt"
	"strings"
)

type Regulation struct {
	Document string `json:"document"`
	Article  string `json:"article"`
	Content  string `json:"content"`
}

type TestCase struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Regulations []Regulation `json:"regulations"`
	Disharmony  string       `json:"disharmony"`
}

const CostarHeader = `
Context:
You are a legal analyst AI specialized in analyzing potential disharmonies or conflicts between regulatory texts in Indonesian law. You are given multiple legal provisions from different legal documents. These provisions might seem aligned but can contain conflicting rules, vague overlaps, or inconsistent exceptions that cause confusion in implementation.

Objective:
Identify and explain any potential legal disharmony, contradiction, or ambiguity between the given legal provisions. Focus on conflicts in meaning, scope, exceptions, or enforcement that may cause practical or legal ambiguity.

Style:
Formal and analytical, using clear and concise legal reasoning.

Tone:
Neutral and professional.

Audience:
Indonesian legal professionals, lawmakers, and policy analysts evaluating coherence across legal regulations.

Response:
Return your answer in Indonesian, plain text, without numbering, bullets, or any other formatting. Provide a clear and concise analysis of the disharmony, including the specific articles or sections that are in conflict, and suggest possible solutions or clarifications to resolve the disharmony.
`

func CostarPromptZeroShot(regulations string) string {
	return fmt.Sprintf("%s\n\nNow, based on the input below, analyze and return your assessment.\n\nInput Regulations:\n%s\n\nEnd of input.", CostarHeader, regulations)
}

func CostarPromptFewShot(regulations string) string {
	var promptBuilder strings.Builder
	promptBuilder.WriteString(CostarHeader)
	promptBuilder.WriteString("\n\nBelow are some examples:")

	testCases, err := LoadTestCases("../data/test_case.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		var exRegText string
		for _, reg := range ex.Regulations {
			exRegText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}
		promptBuilder.WriteString(fmt.Sprintf(
			"\nInput Regulations:\n%s\nExpected Disharmony:\n%s\n\n---\n", exRegText, ex.Disharmony))
	}

	// Add current case (no answer, just input)
	promptBuilder.WriteString("\nNow, analyze the following:\n")
	promptBuilder.WriteString("Input Regulations:\n")
	promptBuilder.WriteString(regulations)
	promptBuilder.WriteString("End of input.\n")

	return promptBuilder.String()
}

func CostarPromptChainOfThought(regulations string) string {
	var cotPromptBuilder strings.Builder
	cotPromptBuilder.WriteString(CostarHeader)
	cotPromptBuilder.WriteString(`

Instructions:
Please follow a step-by-step reasoning process to analyze the input regulations. Begin by identifying each regulation’s main legal point. Then compare them one by one to uncover any conflict, overlap, or ambiguity. Conclude by summarizing the legal disharmony and optionally propose a resolution.

Steps:
1. Identify key legal obligations or permissions in each regulation.
2. Compare similarities and differences in meaning, scope, and exceptions.
3. Detect any conflict, inconsistency, or ambiguity.
4. Provide a reasoned conclusion and suggest clarification if necessary.

Now, analyze the input below step by step.

Input Regulations:
`)

	cotPromptBuilder.WriteString(regulations)
	cotPromptBuilder.WriteString("\nEnd of input.\n")

	return cotPromptBuilder.String()
}

func EvaluatePrompt(method string, regulations string) string {
	var baseInstructions string

	switch method {
	case "zero-shot":
		baseInstructions = CostarPromptZeroShot(regulations)
	case "few-shot":
		baseInstructions = CostarPromptFewShot(regulations)
	case "cot":
		baseInstructions = ""
	case "few-shot-cot":
		baseInstructions = ""
	default:
		baseInstructions = CostarPromptZeroShot(regulations)
	}

	return fmt.Sprintf("%s\n\n%s", baseInstructions, regulations)
}
