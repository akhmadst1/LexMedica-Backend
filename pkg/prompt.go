package pkg

import (
	"fmt"
	"strings"
)

func ZeroShot(regulations string) string {
	return fmt.Sprintf(`Answer in Indonesian, identify the disharmony between this following Indonesian law regulations:
	%s
	End of input.`,
		regulations)
}

func FewShot(regulations string) string {
	var fewShotPromptBuilder strings.Builder
	fewShotPromptBuilder.WriteString("Below are some examples of potential disharmony on Indonesian law regulations:")
	testCases, err := LoadTestCases("./data/test_case.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		var exRegText string
		for _, reg := range ex.Regulations {
			exRegText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}
		fewShotPromptBuilder.WriteString(fmt.Sprintf(
			"\nInput Regulations:\n%s\nExpected Potential Disharmony Analysis Output:\n%s\n\n---\n", exRegText, ex.Disharmony))
	}
	fewShotPromptBuilder.WriteString("\nEnd of examples.\n\n")

	fewShotPromptBuilder.WriteString("Answer in Indonesian, identify the potential disharmony between this following Indonesian law regulations:\n")
	fewShotPromptBuilder.WriteString(regulations)
	fewShotPromptBuilder.WriteString("\nEnd of input.\n")

	return fewShotPromptBuilder.String()
}

func ChainOfThought(regulations string) string {
	var cotPromptBuilder strings.Builder
	cotPromptBuilder.WriteString("Answer in Indonesian, identify the potential disharmony between the following Indonesian law regulations.\n")
	cotPromptBuilder.WriteString(`
	Follow this step by step of reasoning to identify the disharmony:
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
		
		Now, analyze the regulations input below:`)
	cotPromptBuilder.WriteString(regulations)
	cotPromptBuilder.WriteString("\nEnd of input.\n")

	return cotPromptBuilder.String()
}

func FewShotChainOfThought(regulations string) string {
	var fewShotCotPromptBuilder strings.Builder
	fewShotCotPromptBuilder.WriteString("Answer in Indonesian, identify the disharmony between the following Indonesian law regulations.\n")
	fewShotCotPromptBuilder.WriteString(`
	Analyze the following step by step of reasoning:
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
	3. Provide your analysis in Indonesian, plain text, in paragraph, without numbering, bullets, or any other formatting.`)

	fewShotCotPromptBuilder.WriteString("\n\nBelow are some examples:")

	testCases, err := LoadTestCases("./data/test_case.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		var exRegText string
		for _, reg := range ex.Regulations {
			exRegText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}
		fewShotCotPromptBuilder.WriteString(fmt.Sprintf(
			"\nInput Regulations:\n%s\nReasoning:\n%s\nExpected Potential Disharmony Analysis Output:\n%s\n\n---\n", exRegText, ex.Reason, ex.Disharmony))
	}

	fewShotCotPromptBuilder.WriteString("Now, analyze the regulations input below:\n")
	fewShotCotPromptBuilder.WriteString(regulations)
	fewShotCotPromptBuilder.WriteString("\nEnd of input.\n")

	return fewShotCotPromptBuilder.String()
}
