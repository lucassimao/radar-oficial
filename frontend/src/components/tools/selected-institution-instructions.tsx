import { useAssistantToolUI } from "@assistant-ui/react";
import { useInstitution } from "../hooks/useInstitution";
import { ClientOnly } from "../ui/client-only";


 const ToolUI = () => {
  const {institution} = useInstitution();

  useAssistantToolUI({
    toolName: "institution-selected",
    render: () => {

      return (
        <div className="bg-white rounded-lg shadow p-5 max-w-md text-gray-800 space-y-4">
        <h2 className="text-lg font-semibold">📚 Diário selecionado:</h2>
        <p className="text-blue-700 font-bold">{institution?.name}</p>
  
        <p>
          Agora você pode enviar suas perguntas sobre este órgão. As respostas serão
          geradas com base nas publicações mais recentes disponíveis neste Diário Oficial.
        </p>
  
        <p className="text-sm text-gray-600">
          Exemplo: <em>“Houve nomeações recentes nesta semana?”</em>
        </p>
      </div>
      );
    },
  });
  return null;
};


// needed to wrap ToolUI with ClientOnly so that we can use local storage api
export const InstitutionSelectedInstructionsUI = ()=>{
  return <ClientOnly><ToolUI/></ClientOnly>
}  