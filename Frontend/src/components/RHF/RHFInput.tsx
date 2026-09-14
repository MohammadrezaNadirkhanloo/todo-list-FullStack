import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { Info } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Spinner } from "../ui/spinner";
import {
  Controller,
  useFormContext,
  type FieldValues,
  type Path,
} from "react-hook-form";
import { type InputHTMLAttributes } from "react";

interface RHFInputProps<T extends FieldValues> extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  "name"
> {
  name: Path<T>;
  label?: string;
  required?: boolean;
  iconStart?: React.ReactNode;
  iconEnd?: React.ReactNode;
  description?: string;
  loading?: boolean;
  autoComplete?:
    | "off"
    | "on"
    | "current-password"
    | "new-password"
    | "new-username";
}

function RHFInput<T extends FieldValues>({
  name,
  label,
  required = false,
  iconStart,
  iconEnd,
  description,
  loading = false,
  placeholder,
  type = "text",
  disabled,
  autoComplete = "off",
  ...inputProps
}: RHFInputProps<T>) {
  const { control } = useFormContext<T>();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field className="gap-1.5" {...inputProps} data-invalid={fieldState.invalid}>
          {label && (
            <FieldLabel htmlFor={field.name}>
              {label} :{required && <span className="text-destructive">*</span>}
            </FieldLabel>
          )}

          <InputGroup>
            {description && (
              <Popover>
                <PopoverTrigger asChild>
                  <InputGroupAddon>
                    <InputGroupButton variant="secondary" size="icon-xs">
                      <Info />
                    </InputGroupButton>
                  </InputGroupAddon>
                </PopoverTrigger>
                <PopoverContent
                  align="start"
                  className="rounded-xl text-xs border bg-popover p-3"
                >
                  <p>{description}</p>
                </PopoverContent>
              </Popover>
            )}

            <InputGroupInput
              {...field}
              {...inputProps}
              id={field.name}
              type={type}
              value={field.value ?? ""}
              placeholder={placeholder}
              disabled={disabled || loading}
              aria-invalid={fieldState.invalid}
              autoComplete={autoComplete}
            />

            {iconStart && (
              <InputGroupAddon align="inline-start">
                {iconStart}
              </InputGroupAddon>
            )}

            {(iconEnd || loading) && (
              <InputGroupAddon align="inline-end">
                {loading ? <Spinner /> : iconEnd}
              </InputGroupAddon>
            )}
          </InputGroup>

          {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
        </Field>
      )}
    />
  );
}

export default RHFInput;
