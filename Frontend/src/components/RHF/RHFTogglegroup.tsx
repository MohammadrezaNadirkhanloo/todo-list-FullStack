import {
  ToggleGroup,
  ToggleGroupItem,
} from "@/components/ui/toggle-group";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";

interface ToggleOption {
  value: string;
  label: string;
  ariaLabel?: string;
  disabled?: boolean;
}

interface RHFToggleGroupProps<T extends FieldValues> {
  name: Path<T>;
  label?: string;
  required?: boolean;
  options: ToggleOption[];
  type?: "single" | "multiple";
  variant?: "default" | "outline";
  disabled?: boolean;
  className?: string;
}

function RHFToggleGroup<T extends FieldValues>({
  name,
  label,
  required = false,
  options,
  type = "single",
  variant = "outline",
  disabled = false,
  className,
}: RHFToggleGroupProps<T>) {
  const { control } = useFormContext<T>();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field className={className} data-invalid={fieldState.invalid || undefined}>
          {label && (
            <FieldLabel>
              {label} :
              {required && <span className="text-destructive"> *</span>}
            </FieldLabel>
          )}

          {type === "single" ? (
            <ToggleGroup
              type="single"
              variant={variant}
              value={field.value ?? ""}
              onValueChange={(val) => {
                if (!val) return; 
                field.onChange(val);
                field.onBlur();
              }}
              disabled={disabled}
              aria-invalid={fieldState.invalid}
            >
              {options.map((opt) => (
                <ToggleGroupItem
                  key={opt.value}
                  value={opt.value}
                  aria-label={opt.ariaLabel ?? opt.label}
                  disabled={opt.disabled}
                  className="cursor-pointer"
                >
                  {opt.label}
                </ToggleGroupItem>
              ))}
            </ToggleGroup>
          ) : (
            <ToggleGroup
              type="multiple"
              variant={variant}
              value={field.value ?? []}
              onValueChange={(val) => {
                field.onChange(val);
                field.onBlur();
              }}
              disabled={disabled}
              aria-invalid={fieldState.invalid}
            >
              {options.map((opt) => (
                <ToggleGroupItem
                  key={opt.value}
                  value={opt.value}
                  aria-label={opt.ariaLabel ?? opt.label}
                  disabled={opt.disabled}
                >
                  {opt.label}
                </ToggleGroupItem>
              ))}
            </ToggleGroup>
          )}

          {fieldState.invalid && (
            <FieldError errors={[fieldState.error]} />
          )}
        </Field>
      )}
    />
  );
}

export default RHFToggleGroup;