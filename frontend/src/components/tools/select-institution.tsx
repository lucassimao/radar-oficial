
import { useAssistantToolUI, useThreadRuntime } from "@assistant-ui/react";
import { useEffect, useState } from "react";
import { Institution, useInstitution } from "../hooks/useInstitution";
import { ClientOnly } from "../ui/client-only";

const ToolUI = () => {
  const {saveInstitution} = useInstitution();
  
  const [institutions, setInstitutions] = useState<Institution[]>();
  const runtime = useThreadRuntime();

  const onInstitutionSelected = (institution:Institution) =>{
    saveInstitution(institution)

    runtime.append({
      role:'assistant',
      content: [
        {
          type: "tool-call",
          toolName: "institution-selected",
          toolCallId: String(Date.now()),
          argsText: '',
          args: {},
        },
      ],
    });

  }

  useEffect(() => {

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/institutions`)
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Erro");
        }
      })
      .then((body) => {
        const institutions: Institution[] = body.institutions;
        setInstitutions(institutions);
      })
  }, []);


  useAssistantToolUI({
    toolName: "select-institution",
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
                {institutions?.map((option) => (
                  <button
                    key={option.slug}
                    onClick={() => onInstitutionSelected(option)}
                    className="bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition"
                  >
                    {option.name}
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


// needed to wrap ToolUI with ClientOnly so that we can use local storage api
export const SelectInstitutionUI = ()=>{
  return <ClientOnly><ToolUI/></ClientOnly>
}  