import {
  Field,
  FieldContent,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { Switch } from "@/components/ui/switch";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";
import { useTranslation } from "react-i18next";

interface RHFSwitchProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  description?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
}

function RHFSwitch<T extends FieldValues>({
  name,
  label,
  description,
  disabled = false,
  required = false,
  className,
}: RHFSwitchProps<T>) {
  const { control } = useFormContext<T>();
  const { t } = useTranslation("input");
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field
          orientation="horizontal"
          className={className}
          data-invalid={fieldState.invalid || undefined}
        >
          <FieldContent>
            {label && (
              <FieldLabel htmlFor={field.name}>
                {label}
                {required && <span className="text-destructive"> *</span>}
              </FieldLabel>
            )}
            {description &&
              (name === "enabled" ? (
                <FieldDescription>
                  {t("enabled.description", { name: description })}
                </FieldDescription>
              ) : (
                <FieldDescription>{description}</FieldDescription>
              ))}
          </FieldContent>

          <Switch
            id={field.name}
            ref={field.ref}
            checked={!!field.value}
            onCheckedChange={field.onChange}
            onBlur={field.onBlur}
            disabled={disabled}
            aria-invalid={fieldState.invalid}
          />
        </Field>
      )}
    />
  );
}

export default RHFSwitch;
