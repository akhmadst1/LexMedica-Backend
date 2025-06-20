package pkg

import (
	"fmt"
	"strings"
)

const disharmonyPromptHeader = `You are a law expert that can detect potential disharmony in Indonesian health-related regulations.

Potential disharmony refers to a situation where two or more legal provisions appear to regulate a similar subject, but contain inconsistencies, ambiguities, or differences that might create confusion or practical difficulty in implementation.

However, not all differences indicate a potential disharmony. Only identify a potential disharmony if there is a **clear and significant** inconsistency in terms of meaning, scope, timing, authority, or implementation consequences that may result in uncertainty or legal conflict.

You are given several provisions from Indonesian health-related legal documents. Each provision may stand on its own or relate to others. Your task is to review and analyze whether these provisions create a potential legal disharmony, contradiction, or ambiguity.

Be selective: **only identify disharmony when the inconsistency is material and can lead to confusion or conflict in practice.** If the difference is minimal, contextually justified, or not likely to cause implementation issues, state that there is no potential disharmony.

Give your answer in Indonesian language and JSON format like this:

{
	"analysis": "1. <Document Name> Nomor <Document Number> Tahun <Year> Pasal <number of clause> Ayat <number of verse>\n"<the verse text>"\n....continue to all regulations involved>\n\n<string of analysis>,
	"result": <boolean>
}

Format Analysis:
- Always quote the provisions first, then provide analysis in a new paragraph.
- Mention the name and number of the regulation, followed by the relevant article (e.g., "Pasal 31 ayat (2) Peraturan Pemerintah Nomor 61 Tahun 2014").
- Highlight conflicting or critical keywords using **bold** text.
- Clearly describe why there is (or is not) a potential disharmony.
- If no disharmony found, explain briefly why the provisions are considered consistent or not in conflict.

Important Notes:
- Field "analysis": give thoughtful analysis in proper paragraphing. Do not include bullet points or headings.
- Field "result": set to true if disharmony is found based on the analysis, false otherwise.
- Do not add any message outside the JSON output.`

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
	testCases, err := LoadTestCases("data/test_cases_5.json")
	if err != nil {
		return fmt.Sprintf("Error loading test cases: %v", err)
	}

	for _, ex := range testCases {
		// if ex.ID == tcId {
		// 	continue
		// }

		var exRegText string
		for i, reg := range ex.Regulations {
			exRegText += fmt.Sprintf("%d. %s %s:\n %s\n\n", i, reg.Document, reg.Article, reg.Content)
		}
		fewShotPromptBuilder.WriteString(fmt.Sprintf(
			"\n%s\n\n%s\n\n---\n", exRegText, ex.Disharmony))
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

	testCases, err := LoadTestCases("data/test_cases_5.json")
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
