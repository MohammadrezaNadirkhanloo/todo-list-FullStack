import React from "react";
import { Button } from "./ui/button";
import { Spinner } from "./ui/spinner";
import { cn } from "@/lib/utils";

function ButtonComponent({
  isLoadingSend = false,
  lable,
  disabled = false,
  className,
  icon,
  isForm = true,
  ...props
}: {
  isLoadingSend?: boolean;
  lable?: string;
  disabled?: boolean;
  isForm?: boolean;
  icon?: React.ReactNode;
  className?: string;
} & Omit<React.ComponentProps<typeof Button>, "className">) {
  const isDisabled = isLoadingSend || disabled;

  return (
    <Button
      className={cn(
        isForm && "mt-5",
        isDisabled && "cursor-not-allowed",
        className,
      )}
      style={{
        cursor: isDisabled ? "not-allowed" : "pointer",
        ...props.style,
      }}
      disabled={isDisabled}
      {...props}
    >
      {isLoadingSend && <Spinner data-icon="inline-start" />}
      {!isLoadingSend && icon}
      {lable || ""}
    </Button>
  );
}

export default ButtonComponent;
