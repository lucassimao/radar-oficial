import { Institution, useInstitution } from "@/components/hooks/useInstitution";
import { SelectInstitutionUI } from "@/components/tools/select-institution";
import { InstitutionSelectedInstructionsUI } from "@/components/tools/selected-institution-instructions";
import {
  AssistantRuntimeProvider,
  useLocalRuntime,
  ChatModelRunOptions,ChatModelRunResult
} from "@assistant-ui/react";
import { type ReactNode } from "react";


class MyModelAdapter {

  constructor(private institution: Institution|null){  }

  async run({ messages, abortSignal }:ChatModelRunOptions): Promise<ChatModelRunResult> {

    if (!this.institution){
      return {
        content: [
          {
            type: "tool-call",
            toolName: "select-institution",
            toolCallId: String(Date.now()),
            argsText: '',
            args: {},
          },
        ],
      }
    }

    const result = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/chat?i=${this.institution.slug}`, {
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
  const {institution} = useInstitution()
  const runtime = useLocalRuntime(new MyModelAdapter(institution));

  return (
    <AssistantRuntimeProvider  runtime={runtime}>
        <SelectInstitutionUI />
      <InstitutionSelectedInstructionsUI />
      {children}
    </AssistantRuntimeProvider>
  );
}
