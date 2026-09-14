// components/RHF/FormContext.tsx
import { FormProvider as RHFFormProvider } from "react-hook-form";
import { type ReactNode } from "react";
import { type FieldValues, type UseFormReturn, type Control } from "react-hook-form";

interface FormProviderProps<
  T extends FieldValues,
  TContext = any,
  TTransformedValues = T,
> {
  form?: UseFormReturn<T, TContext, TTransformedValues>;
  control?: Control<T, TContext, TTransformedValues>;
  children: ReactNode;
}

export function FormProvider<
  T extends FieldValues,
  TContext = any,
  TTransformedValues = T,
>({
  form,
  control,
  children,
}: FormProviderProps<T, TContext, TTransformedValues>) {
  // اگه form داد از اون استفاده کن، وگرنه یه mock object بساز
  const formValue = form || ({
    control: control!,
  } as UseFormReturn<T, TContext, TTransformedValues>);

  return (
    <RHFFormProvider {...formValue}>
      {children}
    </RHFFormProvider>
  );
}