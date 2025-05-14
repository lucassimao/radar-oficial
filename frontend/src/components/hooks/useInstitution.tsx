import { useLocalStorage } from "@uidotdev/usehooks";
import { Dispatch, SetStateAction } from "react";

type Result =  {
  diarioState: string |null
  diarioStateFullName: string |null
  saveDiarioState: Dispatch<SetStateAction<DiarioState | null>>
}

export const BRAZIL_STATES = {
  "AC": "Acre",
  "AL": "Alagoas",
  "AP": "Amapá",
  "AM": "Amazonas",
  "BA": "Bahia",
  "CE": "Ceará",
  "DF": "Distrito Federal",
  "ES": "Espírito Santo",
  "GO": "Goiás",
  "MA": "Maranhão",
  "MT": "Mato Grosso",
  "MS": "Mato Grosso do Sul",
  "MG": "Minas Gerais",
  "PA": "Pará",
  "PB": "Paraíba",
  "PR": "Paraná",
  "PE": "Pernambuco",
  "PI": "Piauí",
  "RJ": "Rio de Janeiro",
  "RN": "Rio Grande do Norte",
  "RS": "Rio Grande do Sul",
  "RO": "Rondônia",
  "RR": "Roraima",
  "SC": "Santa Catarina",
  "SP": "São Paulo",
  "SE": "Sergipe",
  "TO": "Tocantins"
};

export type DiarioState = keyof typeof BRAZIL_STATES

export function useDiarioState():Result {
  const [diarioState, saveDiarioState] = useLocalStorage<DiarioState | null>(
    "diarioState",
    null
  );

  return {
      diarioState,
      diarioStateFullName: diarioState ? BRAZIL_STATES[diarioState] : null,
    saveDiarioState,
  };
}
