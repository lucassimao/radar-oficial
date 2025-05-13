import { useLocalStorage } from "@uidotdev/usehooks";
import { Dispatch, SetStateAction } from "react";

export type Institution = {
  id: number;
  name: string;
  slug: string;
};

type Result =  {
  institution: Institution |null
  saveInstitution: Dispatch<SetStateAction<Institution | null>>
}

export function useInstitution():Result {
  const [institution, saveInstitution] = useLocalStorage<Institution | null>(
    "institution",
    null
  );

  return {
    institution,
    saveInstitution,
  };
}
