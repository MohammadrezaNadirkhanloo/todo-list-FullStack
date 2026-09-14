import { FormProvider } from "@/components/RHF/FormContext";
import RHFInput from "@/components/RHF/RHFInput";
import RHFPasswordInput from "@/components/RHF/RHFPasswordInput";
import { Lock, User } from "lucide-react";
import { useCreateUsers } from "../api/postLoginApi";
import { useLoginForm } from "./useLoginForm";
import type { LoginFormData } from "./auth.schema";
import ButtonComponent from "@/components/Button-custom";
import { AuthDivider } from "@/components/auth-divider";
import { Separator } from "@/components/ui/separator";

export function SingInForm() {
  const { LoginUsers, isLoadingSend } = useCreateUsers();
  const form = useLoginForm();

  function onSubmit(data: LoginFormData) {
    LoginUsers(data);
  }

  return (
    <FormProvider form={form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="grid grid-cols-1 w-full max-w-sm mx-auto px-5 gap-3"
      >
        <RHFInput
          name="username"
          placeholder={"example@gmail.com"}
          iconStart={<User className="text-muted-foreground" />}
          loading={isLoadingSend}
          
        />
        <RHFPasswordInput
          name="password"
          placeholder={"password"}
          iconStart={<Lock className="h-4 w-4 text-muted-foreground" />}
          autoComplete="current-password"
          loading={isLoadingSend}
        />
        <Separator />
        <ButtonComponent
          lable="Login"
          isLoadingSend={isLoadingSend}
          type="submit"
          className="mt-0 sm:min-w-28 w-full"
          // icon={isEdit ? <Save /> : iconCreate}
        />
        {/* <ButtonComponent
          lable="login"
          isLoadingSend={isLoadingSend}
          type="submit"
        /> */}
      </form>
    </FormProvider>
  );
}
