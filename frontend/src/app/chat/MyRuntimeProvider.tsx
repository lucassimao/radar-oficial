import { useDiarioState } from "@/components/hooks/useInstitution";
import { SELECT_DIARIO_STATE_TOOL_NAME, SelectStateUI } from "@/components/tools/select-state";
import { StateSelectedInstructionsUI } from "@/components/tools/selected-state-instructions";
import {
  AssistantRuntimeProvider,
  ChatModelRunOptions, ChatModelRunResult,
  useLocalRuntime
} from "@assistant-ui/react";
import { type ReactNode } from "react";


class MyModelAdapter {

  constructor(private diarioState: string|null){  }

  async run({ messages, abortSignal }:ChatModelRunOptions): Promise<ChatModelRunResult> {

    if (!this.diarioState){
      return {
        content: [
          {
            type: "tool-call",
            toolName: SELECT_DIARIO_STATE_TOOL_NAME,
            toolCallId: String(Date.now()),
            argsText: '',
            args: {},
          },
        ],
      }
    }

    const result = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/chat?state=${this.diarioState}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      // forward the messages in the chat to the API
      body: JSON.stringify({
        messages,
      }),
      // if the user hits the "cancel" button or escape keyboard key, cancel the request
      signal: abortSignal,
    });

    const data = await result.json();
    return {
      content: [
        {
          type: "text",
          text: data.text,
        },
      ],
    }
  }
};

export function MyRuntimeProvider({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  const {diarioState} = useDiarioState()
  const runtime = useLocalRuntime(new MyModelAdapter(diarioState));

  return (
    <AssistantRuntimeProvider  runtime={runtime}>
        <SelectStateUI />
      <StateSelectedInstructionsUI />
      {children}
    </AssistantRuntimeProvider>
  );
}
