import { useAssistantToolUI } from "@assistant-ui/react";
import { useDiarioState } from "../hooks/useInstitution";


export const DIARIO_STATE_SELECTED_TOOL_NAME= 'diario-state-selected'

export const StateSelectedInstructionsUI = () => {
  const {diarioStateFullName} = useDiarioState();

  useAssistantToolUI({
    toolName: DIARIO_STATE_SELECTED_TOOL_NAME,
    render: () => {

      return (
        <div className="bg-white rounded-lg shadow p-5 max-w-md text-gray-800 space-y-4">
        <h2 className="text-lg font-semibold">📚 Diário selecionado!</h2>
        {/* <p className="text-blue-700 font-bold">{institution?.name}</p> */}
  
        <p>
          Agora você pode enviar suas perguntas. As respostas serão
          geradas com base em todas as publicações oficiais mais recentes do estado 
          de {diarioStateFullName} .
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

