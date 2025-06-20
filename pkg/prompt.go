package pkg

import (
	"fmt"
	"strings"
)

const disharmonyPromptHeader = `You are a law expert that can detect potential disharmony problem in Indonesian law documents.

Potential disharmony is a condition in which two or more regulations address similar subject matter but are inconsistent in their technical specifications.

Fundamentally, this creates conflicts between regulations and leads to setbacks either horizontally (across sectors or institutions) or vertically (between hierarchical levels of law).

You are given multiple legal provisions from Indonesian legal documents. These provisions might seem aligned but can contain conflicting rules, vague overlaps, or inconsistent exceptions that cause confusion in implementation.

Your task is to identify and explain any potential legal disharmony, contradiction, or ambiguity between the given legal provisions.
You don't have to be strict, most of the time the regulations is not really potentially disharmony even if it seems to be, i'll give you some examples later on, if it's not really close with the example then it's not potentially disharmony.
Focus on conflicts in meaning, scope, exceptions, or enforcement that may cause practical or legal ambiguity.

Give your answer in Indonesian language and JSON format like this:

{
"result": <boolean>,
"analysis": <string>
}

Important Notes:
- Not all case has potential disharmony, state it clearly if there is not any potential disharmony found.
- Field "result" is whether you found potential disharmony or not, if disharmony found set to true, if not set to false.
- Field "analysis" is the analysis text of the potential disharmony.
- State the regulations with quotation mark, and bold the keywords that make those potentially disharmony.
- Make new paragraph if the analysis contain more than one case.
- Do not add any other message outside the JSON format.`

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
	testCases, err := LoadTestCases("data/test_cases_7.json")
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
1. Identify key norms from each regulation.
2. Compare definitions, scope, responsibilities, and exceptions.
3. Look for entity and description comparisons:
	- Entities: people, organizations, time, locations.
	- Descriptions: what is the subject, object, and action in each regulation.
4. Look for overlaps or mismatches between entities and descriptions.
5. If any, extract and consider the comparison between following:
	- Nominal values (currency or quantities) mentioned in the provisions.
	- Positions or official titles (e.g., Menteri, Kepala Dinas, etc.).
	- References to other articles, clauses, or laws (e.g., "Pasal 3 ayat (2)").
6. If any, evaluate and compare the number of requirements between the related provisions.
7. Conclude whether there is any potential disharmony and explain the reasoning in detail.
8. Do not make any suggestion, just state the comparison.

Now, analyze the regulations input below:`)
	cotPromptBuilder.WriteString(regulations)
	cotPromptBuilder.WriteString("\nEnd of input.")

	return cotPromptBuilder.String()
}

func FewShotChainOfThought(regulations string) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString(disharmonyPromptHeader)
	promptBuilder.WriteString("\nBelow are some examples:")

	testCases, err := LoadTestCases("data/test_cases_7.json")
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

	promptBuilder.WriteString(`Now, analyze the regulations input below:`)
	promptBuilder.WriteString(regulations)
	promptBuilder.WriteString("\nEnd of input.")

	return promptBuilder.String()
}
