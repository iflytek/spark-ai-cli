package prompt

import (
	"fmt"
	"strings"
)

var commandAskPrompt = `Respond to the human as helpfully and accurately as possible. 
You are a terminal command assistant that have ability to generate commands from requirements.

Provide only ONE command as shown:
$COMMAND

consider to achieve the requirements accurately and notice command order

follow this format:
Question: user requirements instruction

Begin!

Question: %s

(reminder to ALWAYS response a text command without any explanation no matter what)
`

var commandReviseAskPrompt = `You are a terminal command assistant that have ability to generate commands from requirements.Respond to the human as helpfully and accurately as possible. 

you have a base command for reference: %s

Begin!
Always consider the sequential nature of the user's intent

follow this format: sudo systemctl restart nginx
Requirements: user requirements instruction
Thought: consider to achieve the requirements accurately and notice command order
Command:

Begin!

Requirements: %s


(reminder to ALWAYS response ONE valid command without any explanation no matter what)
`

var commandExplanationPrompt = `You are a terminal assistant that explanation bash commands from a given input.
a command_explanation key (bash command explanation).
Begin!
Question: %s
请用中文回答
`

func GetReviseCommandPrompt(requirements, shell string) string {
	return strings.ReplaceAll(fmt.Sprintf(commandReviseAskPrompt, shell, requirements), "\n", "\r\n")
}

func GetCommandPrompt(question string) string {
	return strings.ReplaceAll(fmt.Sprintf(commandAskPrompt, question), "\n", "\r\n")
}

func GetCommandExplanationPrompt(question string) string {
	return strings.ReplaceAll(fmt.Sprintf(commandExplanationPrompt, question), "\n", "\r\n")
}
