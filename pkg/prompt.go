package pkg

import (
	"fmt"
	"strings"
)

const disharmonyPromptHeader = `You are a law expert that can detect potential disharmony problem in law documents.
	Disharmony is a condition in which two or more regulations address similar subject matter but are inconsistent in their technical specifications.
	Fundamentally, this creates conflicts between regulations and leads to setbacks either horizontally (across sectors or institutions) or vertically (between hierarchical levels of law).
	You are given multiple legal provisions from Indonesian legal documents. These provisions might seem aligned but can contain conflicting rules, vague overlaps, or inconsistent exceptions that cause confusion in implementation.
	Your task is to identify and explain any potential legal disharmony, contradiction, or ambiguity between the given legal provisions.
	Focus on conflicts in meaning, scope, exceptions, or enforcement that may cause practical or legal ambiguity.

	Give your answer in Indonesian language and JSON format like this:
	
	{
	"result": <boolean>,
	"analysis": <string>
	}

	- Field "result" is whether you found potential disharmony or not, if disharmony found set to true, if not set to false.
	- Field "analysis" the result of the analysis and must be a brief summary of your step-by-step reasoning.
	- Do not add any other message outside the JSON format.
	`

func ZeroShot(regulations string) string {
	return fmt.Sprintf(`%s
	
	Answer in Indonesian, identify the disharmony between this following Indonesian law regulations:
	
	%s

	End of input.`,
		disharmonyPromptHeader, regulations)
}

func FewShot(regulations string) string {
	var fewShotPromptBuilder strings.Builder
	fewShotPromptBuilder.WriteString(disharmonyPromptHeader)
	fewShotPromptBuilder.WriteString("\nBelow are some examples of potential disharmony on Indonesian law regulations:")
	testCases, err := LoadTestCases("../data/test_case.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		// if ex.ID == tcId {
		// 	continue
		// }

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
	fewShotPromptBuilder.WriteString("\nEnd of input.")

	return fewShotPromptBuilder.String()
}

func ChainOfThought(regulations string) string {
	var cotPromptBuilder strings.Builder
	cotPromptBuilder.WriteString(disharmonyPromptHeader)
	cotPromptBuilder.WriteString(`
		Analyze step-by-step:
		- Identify key norms from each regulation.
		- Compare definitions, scope, responsibilities, and exceptions.
		- Find potential contradictions or overlaps.
		- Conclude whether there's disharmony and explain why.
		- Optionally suggest a resolution if needed.
		
		Now, analyze the regulations input below:`)
	cotPromptBuilder.WriteString(regulations)
	cotPromptBuilder.WriteString("\nEnd of input.")

	return cotPromptBuilder.String()
}

func FewShotChainOfThought(regulations string) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString(disharmonyPromptHeader)
	promptBuilder.WriteString(`
	Below are some examples:`)

	testCases, err := LoadTestCases("../data/test_case.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		// if ex.ID == tcId {
		// 	continue
		// }

		var exRegText string
		for _, reg := range ex.Regulations {
			exRegText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}

		promptBuilder.WriteString(fmt.Sprintf(
			"\nInput Regulations:\n%s"+
				"Step-by-step Reasoning:\n%s\n"+
				"Final Output:\n{\n  \"result\": true,\n  \"analysis\": \"%s\"\n}\n"+
				"\n---", exRegText, ex.Reason, strings.ReplaceAll(ex.Disharmony, "\"", "'")))
	}

	promptBuilder.WriteString(`
	Now, analyze the regulations input below:
	`)
	promptBuilder.WriteString(regulations)
	promptBuilder.WriteString("\nEnd of input.")

	return promptBuilder.String()
}
