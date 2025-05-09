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
You are a legal analyst AI specialized in analyzing potential disharmonies or conflicts between regulatory texts in Indonesian law. Regulatory disharmony refers to a condition in which two or more regulations address similar subject matter but are inconsistent in their technical specifications. Fundamentally, this creates conflicts between regulations and leads to setbacks either horizontally (across sectors or institutions) or vertically (between hierarchical levels of law). You are given multiple legal provisions from Indonesian legal documents. These provisions might seem aligned but can contain conflicting rules, vague overlaps, or inconsistent exceptions that cause confusion in implementation. 

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

	testCases, err := LoadTestCases("./data/test_case.json")
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
Now, analyze the following step by step of reasoning:
1. Identify the main legal norms, obligations, or permissions from each regulation, including relevant articles or clauses.
2. Compare these legal norms side by side in terms of:
   - Definition or terminology
   - Legal scope and applicability
   - Authority/responsibility assigned
   - Exceptions or special conditions
3. Analyze potential disharmony:
   - Is there a direct contradiction?
   - Is there overlapping jurisdiction or authority?
   - Are there ambiguities that can lead to multiple interpretations?
4. Conclude with a summary of your findings:
   - Clearly state if disharmony exists or not.
   - If disharmony exists, briefly explain the legal or practical implication.
5. Recommend possible resolutions if appropriate:
   - Suggest harmonization steps (e.g., revision, repeal, new regulation, clarification).
   - Mention which regulation may take precedence (if applicable).

Important Notes:
1. Be neutral and objective.
2. Not all case will have disharmony, so if you find no disharmony, please state that clearly.
3. Provide your analysis in Indonesian, plain text, in paragraph, without numbering, bullets, or any other formatting.

Now, analyze the regulations input below:
`)
	cotPromptBuilder.WriteString(regulations)
	cotPromptBuilder.WriteString("\nEnd of input.\n")

	return cotPromptBuilder.String()
}

func CostarPromptFewShotChainOfThought(regulations string) string {
	var promptBuilder strings.Builder
	promptBuilder.WriteString(CostarHeader)
	promptBuilder.WriteString(`
Now, analyze the following step by step of reasoning:
1. Identify the main legal norms, obligations, or permissions from each regulation, including relevant articles or clauses.
2. Compare these legal norms side by side in terms of:
   - Definition or terminology
   - Legal scope and applicability
   - Authority/responsibility assigned
   - Exceptions or special conditions
3. Analyze potential disharmony:
   - Is there a direct contradiction?
   - Is there overlapping jurisdiction or authority?
   - Are there ambiguities that can lead to multiple interpretations?
4. Conclude with a summary of your findings:
   - Clearly state if disharmony exists or not.
   - If disharmony exists, briefly explain the legal or practical implication.
5. Recommend possible resolutions if appropriate:
   - Suggest harmonization steps (e.g., revision, repeal, new regulation, clarification).
   - Mention which regulation may take precedence (if applicable).

Important Notes:
1. Be neutral and objective.
2. Not all case will have disharmony, so if you find no disharmony, please state that clearly.
3. Provide your analysis in Indonesian, plain text, in paragraph, without numbering, bullets, or any other formatting.
`)
	promptBuilder.WriteString("\n\nBelow are some examples:")

	testCases, err := LoadTestCases("./data/test_case.json")
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

	promptBuilder.WriteString("Now, analyze the regulations input below:\n")
	promptBuilder.WriteString("Input Regulations:\n")
	promptBuilder.WriteString(regulations)
	promptBuilder.WriteString("End of input.\n")

	return promptBuilder.String()
}

func EvaluatePrompt(method string, regulations string) string {
	var baseInstructions string

	switch method {
	case "zero-shot":
		baseInstructions = CostarPromptZeroShot(regulations)
	case "few-shot":
		baseInstructions = CostarPromptFewShot(regulations)
	case "cot":
		baseInstructions = CostarPromptChainOfThought(regulations)
	case "few-shot-cot":
		baseInstructions = CostarPromptFewShotChainOfThought(regulations)
	default:
		baseInstructions = CostarPromptZeroShot(regulations)
	}

	return fmt.Sprintf("%s\n\n%s", baseInstructions, regulations)
}
