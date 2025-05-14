
import { useAssistantToolUI, useThreadListItemRuntime, useThreadRuntime } from "@assistant-ui/react";
import { useEffect, useState } from "react";
import {  BRAZIL_STATES, DiarioState, useDiarioState } from "../hooks/useInstitution";
import { DIARIO_STATE_SELECTED_TOOL_NAME } from "./selected-state-instructions";

export const SELECT_DIARIO_STATE_TOOL_NAME= 'select-diario-state'


export const SelectStateUI = () => {
  const {saveDiarioState} = useDiarioState();
  const threadListItemRuntime = useThreadListItemRuntime();
  const [diarioState, setDiarioStates] = useState<DiarioState[]>();
  const runtime = useThreadRuntime();

  const onDiarioStateSelected = (diarioState:DiarioState) =>{
    saveDiarioState(diarioState)

    threadListItemRuntime.rename(BRAZIL_STATES[diarioState])
    runtime.append({
      role:'assistant',
      content: [
        {
          type: "tool-call",
          toolName: DIARIO_STATE_SELECTED_TOOL_NAME,
          toolCallId: String(Date.now()),
          argsText: '',
          args: {},
        },
      ],
    });

  }

  useEffect(() => {

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/states`)
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Erro");
        }
      })
      .then((body) => {
        const states: DiarioState[] = body.states;
        setDiarioStates(states);
      })
  }, []);


  useAssistantToolUI({
    toolName: SELECT_DIARIO_STATE_TOOL_NAME,
    render: () => {

      return (
        <div className="flex w-full max-w-[var(--thread-max-width)] flex-grow flex-col">
          <div className="flex w-full flex-grow flex-col items-center justify-center">
            <div className="bg-white shadow-lg rounded-xl p-6">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">
                📰 Bem-vindo ao Radar Oficial!
              </h2>
              <p className="text-gray-600 mb-6">
                Escolha o Diário que você deseja consultar:
              </p>
              <div className="flex flex-col gap-3">
                {diarioState?.map((name) => (
                  <button
                    key={name}
                    onClick={() => onDiarioStateSelected(name)}
                    className="bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition"
                  >
                    {BRAZIL_STATES[name]}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      );
    },
  });
  return null;
};

