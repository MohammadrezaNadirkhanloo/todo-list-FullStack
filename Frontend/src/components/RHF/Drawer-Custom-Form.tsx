import { useLanguage } from "@/context/LanguageContext";
import {
  cloneElement,
  isValidElement,
  type ReactElement,
  type ReactNode,
} from "react";
import { useTranslation } from "react-i18next";
import { Button } from "../ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "../ui/sheet";

function DrawerCustomForm({
  textButton,
  title,
  description = true,
  children,
  iconButton,
  sizeIcon = "icon",
  variant = "outline",
  size = "xl",
  className,
  id,
  open,
  onOpenChange,
}: {
  textButton?: string;
  title?: string;
  description?: boolean;
  children: ReactNode;
  iconButton?: ReactNode;
  sizeIcon?: "icon" | "icon-sm";
  variant?: "outline" | "secondary" | "ghost" | "default";
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
  id?: string | number;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}) {
  const { direction } = useLanguage();
  const { t } = useTranslation("common");
  const isControlled = open !== undefined;

  const content =
    id !== undefined && isValidElement(children)
      ? cloneElement(children as ReactElement<{ id?: string | number }>, {
          id,
        })
      : children;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      {!isControlled && (
        <SheetTrigger asChild>
          {iconButton && textButton ? (
            <Button variant={variant} size="sm" className={className}>
              {iconButton}
              <span className={"hidden sm:inline-block"}>
                {t(`sheet.actions.${textButton}`)}
              </span>
            </Button>
          ) : iconButton ? (
            <Button variant={variant} size={sizeIcon}>
              {iconButton}
            </Button>
          ) : (
            <Button variant={variant}>{t(`sheet.actions.${textButton}`)}</Button>
          )}
        </SheetTrigger>
      )}
      <SheetContent size={size} dir={direction}>
        {(title || description) && (
          <SheetHeader>
            {title && <SheetTitle>{t(`sheet.title.${title}`)}</SheetTitle>}
            {description && (
              <SheetDescription>
                {t(`sheet.description.${title}`)}
              </SheetDescription>
            )}
          </SheetHeader>
        )}
        {content}
      </SheetContent>
    </Sheet>
  );
}

export default DrawerCustomForm;
