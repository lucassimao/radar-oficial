"use client";
 
import { SelectInstitutionUI } from "@/components/tools/select-institution";
import { InstitutionSelectedInstructionsUI } from "@/components/tools/selected-institution-instructions";
import {
  AssistantRuntimeProvider,
  useLocalRuntime,
  type ChatModelAdapter
} from "@assistant-ui/react";
import { Suspense, type ReactNode } from "react";
 
const MyModelAdapter: ChatModelAdapter = {
  async run({ messages, abortSignal }) {
    
    const result = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/chat`, {
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
    };
  },
};

 
export function MyRuntimeProvider({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  const runtime = useLocalRuntime(MyModelAdapter);

  return (
    <AssistantRuntimeProvider  runtime={runtime}>
      <Suspense fallback={<div/>}>
        <SelectInstitutionUI/>
      </Suspense>
      <InstitutionSelectedInstructionsUI/>
      {children}
    </AssistantRuntimeProvider>
  );
}
