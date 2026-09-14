// components/RHF/RHFPasswordInput.tsx
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { Info, Eye, EyeOff } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Spinner } from "../ui/spinner";
import { Controller, useFormContext, type FieldValues, type Path } from "react-hook-form";
import { type InputHTMLAttributes, useState } from "react";

interface RHFPasswordInputProps<T extends FieldValues> extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  "name" | "type"
> {
  name: Path<T>;
  label?: string;
  required?: boolean;
  iconStart?: React.ReactNode;
  description?: string;
  loading?: boolean;
  showToggle?: boolean;
  autoComplete?: "off" | "current-password" | "new-password";
}

function RHFPasswordInput<T extends FieldValues>({
  name,
  label,
  required = false,
  iconStart,
  description,
  loading = false,
  placeholder = "••••••••",
  disabled,
  showToggle = true,
  autoComplete = "off",
  ...inputProps
}: RHFPasswordInputProps<T>) {
  const { control } = useFormContext<T>();
  const [showPassword, setShowPassword] = useState(false);

  const togglePasswordVisibility = () => {
    setShowPassword((prev) => !prev);
  };

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Field className="gap-1.5" data-invalid={fieldState.invalid}>
          <FieldLabel htmlFor={field.name}>
            {label}
            {required && <span className="text-destructive ml-1">*</span>}
          </FieldLabel>

          <InputGroup>
            {/* دکمه توضیحات */}
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

            {/* Input پسورد */}
            <InputGroupInput
              {...field}
              {...inputProps}
              id={field.name}
              type={showPassword ? "text" : "password"}
              placeholder={placeholder}
              disabled={disabled || loading}
              aria-invalid={fieldState.invalid}
              autoComplete={autoComplete}
            />

            {/* آیکون شروع */}
            {iconStart && (
              <InputGroupAddon align="inline-start">
                {iconStart}
              </InputGroupAddon>
            )}

            {/* دکمه نمایش/مخفی پسورد یا Spinner */}
            <InputGroupAddon align="inline-end">
              {loading ? (
                <Spinner />
              ) : showToggle ? (
                <InputGroupButton
                  type="button"
                  variant="ghost"
                  size="icon-xs"
                  onClick={togglePasswordVisibility}
                  disabled={disabled}
                  aria-label={showPassword ? "مخفی کردن پسورد" : "نمایش پسورد"}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4 text-muted-foreground" />
                  ) : (
                    <Eye className="h-4 w-4 text-muted-foreground" />
                  )}
                </InputGroupButton>
              ) : null}
            </InputGroupAddon>
          </InputGroup>

          {/* نمایش خطا */}
          {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
        </Field>
      )}
    />
  );
}

export default RHFPasswordInput;
