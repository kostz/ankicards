package internal

const PromptExtractVerbsAndTranslate = `
<LLMGuidelines>
    <preamble>You have scanned pages from the German schoolbook vocabulary containing german irregular verbs in three columns:
infinitive, present tense and past tense.
    </preamble>
    <InputData>
        <Description>Scan of a german irregular verbs vocabulary page having three columns</Description>
		<VerbColumn>
			<Column>
				<name>infinitive</name>
				<description>verb in infinitive</description>
			</Column>
			<Column>
				<name>present tense</name>
				<description>verb in present tense</description>
			</Column>
			<Column>
				<name>past tense</name>
				<description>verb in past tense, prateritum</description>
			</Column>
		</VerbColumn>
    </InputData>
    <ResponseGuidelines>
        <Task>
            <description>Parse given page extracting values from columns into a json</description>
            <Instructions>
                <Instruction>From the supplied page scan extract values of all three columns into a json structure.
                </Instruction>
				<Instruction>Placing infinitive into "infinitive", present tense into the "present", past tense into "past" tense fields of JSON element
                </Instruction>
            </Instructions>
        </Task>
        <Task>
            <description>Add translation of each verb into Russian and English</description>
            <Instructions>
                <Instruction>For every extracted verb add translation into Russian and English
                </Instruction>
				<Instruction>Place translation into the nested JSON element "translation" for the verb record having fields "ru" for Russian translation and "en" for English translation.
                </Instruction>
            </Instructions>
        </Task>

		<Restrictions>
			  <Instruction>Refrain from surrounding the response in a code block.</Instruction>
			  <Instruction>Refrain from escaping special characters.</Instruction>
			  <Instruction>Refrain from using Markdown or any special Text formating</Instruction>
			  <Instruction>Refrain from adding LLM greetings, closings, or any other non-technical information.
			  </Instruction>
			  <Instruction>Refrain from giving performance recommendations, architecture changes or implementing fallbacks
			  </Instruction>
			  <Instruction>Refrain from suggesting changes to the studying process</Instruction>
			</Restrictions>
    </ResponseGuidelines>
</LLMGuidelines>
`

const PromptAddExampleSentences = `
<LLMGuidelines>
    <preamble>You have JSON structure with a German verb having it's infinitive form, present tense, past tense and translation into Russian and English. '
    </preamble>
    <InputData>
        <Description>JSON structure with German verb</Description>
		<VerbElement>
			<Element>
				<name>infinitive</name>
				<description>verb in infinitive</description>
			</Element>
			<Element>
				<name>present</name>
				<description>verb in present tense</description>
			</Element>
			<Element>
				<name>past</name>
				<description>verb in past tense, prateritum</description>
			</Element>
			<Element>
				<name>translation</name>
				<description>translation into Russian and English languages</description>
			</Element>
		</VerbElement>
    </InputData>
    <ResponseGuidelines>
		<Task>
            <description>Add examples for both present and past tense</description>
            <Instructions>
                <Instruction>For both present and past tense of the verb prepare the example sentence illustrating the usage of a verb. 
                </Instruction>
                <Instruction>Place examples into the nested JSON element "examples" as an array of JSON element containing the example sentence in element "sentence", don't use any interim layer of structure.
                </Instruction>
                <Instruction>For each example sentence make a translation into Russian and English.
                </Instruction>
                <Instruction>Place translation into the nested JSON element "translation" having fields "ru" for Russian translation and "en" for English translation.
                </Instruction>
            </Instructions>
        </Task>
        <Format>
			JSON Structure
        </Format>
		<Restrictions>
			  <Instruction>Refrain from surrounding the response in a code block.</Instruction>
			  <Instruction>Refrain from escaping special characters.</Instruction>
			  <Instruction>Refrain from using Markdown or any special Text formating</Instruction>
			  <Instruction>Refrain from adding LLM greetings, closings, or any other non-technical information.
			  </Instruction>
			  <Instruction>Refrain from giving performance recommendations, architecture changes or implementing fallbacks
			  </Instruction>
			  <Instruction>Refrain from suggesting changes to the studying process</Instruction>
			</Restrictions>
    </ResponseGuidelines>
</LLMGuidelines>
`
