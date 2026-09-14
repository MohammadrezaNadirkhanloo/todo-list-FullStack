import { redirect, type LoaderFunctionArgs } from "react-router";
import { subject as toSubject } from "@casl/ability";
import type { Actions, Subjects } from "@/types";
import { useSessionStore } from "@/zustand/permission.store";

/** نگهبان ساده — دسترسی به یک نوع صفحه */
export function requireAbility(action: Actions, subj: Subjects) {
  return async (): Promise<null> => {
    await useSessionStore.getState().ensureSession();
    const { ability } = useSessionStore.getState();

    if (ability.cannot(action, subj)) {
      throw redirect("/403");
    }
    return null;
  };
}


export function requireRecordAbility<T extends object>(
  action: Actions,
  type: Subjects,
  getRecord: (params: LoaderFunctionArgs["params"]) => T,
) {
  return async ({ params }: LoaderFunctionArgs): Promise<null> => {
    await useSessionStore.getState().ensureSession();
    const { ability } = useSessionStore.getState();

    const tagged = toSubject(type, getRecord(params));
    
    if (ability.cannot(action, tagged as unknown as Subjects)) {
      throw redirect("/403");
    }
    return null;
  };
}