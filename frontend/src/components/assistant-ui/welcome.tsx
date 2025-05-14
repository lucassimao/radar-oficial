import { useThreadRuntime } from "@assistant-ui/react";
import React from "react";
import { SELECT_DIARIO_STATE_TOOL_NAME } from "../tools/select-state";

export const WelcomeChatbot: React.FC<{}> = () => {
  const runtime = useThreadRuntime();

  const onSubscribe = () => {
    // TODO subscribe flow
  };

  const onClickListDiarios = () => {
    runtime.append({
      role:'assistant',
      content: [
        {
          type: "tool-call",
          toolName: SELECT_DIARIO_STATE_TOOL_NAME,
          toolCallId: String(Date.now()),
          argsText: '',
          args: {},
        },
      ],
    });
  };

  return (
    <div className="bg-white rounded-lg shadow p-5 max-w-md text-gray-800 space-y-4">
      <h2 className="text-xl font-bold">👋 Bem-vindo ao Radar Oficial!</h2>

      <p>
        O <strong>Radar Oficial</strong> foi criado para ajudar você a encontrar
        com facilidade informações publicadas em{" "}
        <strong>Diários Oficiais</strong> de todo o <strong>Brasil</strong>.
      </p>

      <p>
        Basta dizer o que procura — e nosso assistente te ajuda a encontrar o
        conteúdo nos documentos oficiais.
      </p>

      <div className="border-t pt-4">
        <h3 className="font-semibold mb-2">📦 Planos disponíveis</h3>
        <ul className="space-y-1 text-sm">
          <li>
            ✅ <strong>Gratuito</strong>: até 5 consultas por mês
          </li>
          <li>
            ✅ <strong>Básico</strong>: até 5 consultas por dia (R$ 19,90/mês)
          </li>
          <li>
            ✅ <strong>Profissional</strong>: Consultas ilimitadas (R$ 49,90/mês)
          </li>
        </ul>
      </div>

      <div className="flex flex-col md:flex-row gap-3 pt-4">
        <button
          onClick={onSubscribe}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition"
        >
          Assinar um plano
        </button>
        <button
          onClick={onClickListDiarios}
          className="bg-gray-200 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-300 transition"
        >
          Listar Diários disponíveis
        </button>
      </div>
    </div>
  );
};
