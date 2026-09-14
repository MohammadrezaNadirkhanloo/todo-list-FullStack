import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/components/ui/field";
import { Textarea } from "@/components/ui/textarea";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";
import { type TextareaHTMLAttributes } from "react";

interface RHFTextareaProps<T extends FieldValues>
  extends Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, "name"> {
  name: Path<T>;
  label?: string;
  description?: string;
  required?: boolean;
  className?: string;
}

function RHFTextarea<T extends FieldValues>({
  name,
  label,
  description,
  required = false,
  placeholder,
  disabled,
  rows,
  className,
  ...textareaProps
}: RHFTextareaProps<T>) {
  const { control } = useFormContext<T>();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field
          // className={className}
          data-invalid={fieldState.invalid || undefined}
        >
          {label && (
            <FieldLabel htmlFor={field.name}>
              {label}
              {required && <span className="text-destructive"> *</span>}
            </FieldLabel>
          )}

          <Textarea
            {...field}
            {...textareaProps}
            className={className}
            id={field.name}
            placeholder={placeholder}
            disabled={disabled}
            rows={rows}
            aria-invalid={fieldState.invalid}
            value={field.value ?? ""}
          />

          {description && !fieldState.invalid && (
            <FieldDescription className="text-xs">{description}</FieldDescription>
          )}

          {fieldState.invalid && (
            <FieldError errors={[fieldState.error]} />
          )}
        </Field>
      )}
    />
  );
}

export default RHFTextarea;