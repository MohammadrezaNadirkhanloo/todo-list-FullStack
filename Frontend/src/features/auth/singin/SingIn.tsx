import ButtonComponent from "@/components/Button-custom";
import { DecorIcon } from "@/components/decor-icon";
import { cn } from "@/lib/utils";
import { ArrowLeft } from "lucide-react";
import { Link, useNavigate } from "react-router";
import { SingInForm } from "./form/FormLogin";

function SingIn() {
  const navigate = useNavigate();

  return (
    <div className="relative flex h-screen w-full items-center justify-center overflow-hidden px-6 md:px-8">
      <div
        className={cn(
          "relative flex w-full max-w-sm flex-col justify-between p-6 md:p-8",
          "dark:bg-[radial-gradient(50%_80%_at_20%_0%,--theme(--color-foreground/.1),transparent)]",
        )}
      >
        {/* Border decorations */}
        <div className="absolute -inset-y-6 -left-px w-px bg-border" />
        <div className="absolute -inset-y-6 -right-px w-px bg-border" />
        <div className="absolute -inset-x-6 -top-px h-px bg-border" />
        <div className="absolute -inset-x-6 -bottom-px h-px bg-border" />

        <DecorIcon position="top-left" />
        <DecorIcon position="bottom-right" />

        <div className="w-full max-w-sm animate-in space-y-8">
          {/* Header */}
          <div className="flex flex-col space-y-1">
            <div className="flex items-center gap-2.5">
              <ButtonComponent
                icon={<ArrowLeft />}
                size={"icon"}
                className="mt-0"
                onClick={() => navigate("/")}
              />
              <h1 className="text-2xl font-bold tracking-wide">Welcome back</h1>
            </div>

            <p className="text-base text-muted-foreground">
              Sign in to stay organized and keep track of your tasks.
            </p>
          </div>
          <div className="space-y-4">
            <SingInForm />

            {/* <AuthDivider>OR</AuthDivider> */}

            {/* Social Login */}
            {/* <div className="grid grid-cols-2 gap-2">
              <Button variant="outline" className="w-full">
                <GithubIcon className="mr-2 size-4" />
                Github
              </Button>

              <Button variant="outline" className="w-full">
                <GoogleIcon className="mr-2 size-4" />
                Google
              </Button>
            </div> */}
          </div>
          {/* Footer */}
          <div className="space-y-4">
            {/* Sign Up */}
            <p className="text-center text-sm text-muted-foreground">
              Don&apos;t have an account?{" "}
              <Link
                to="/sign-up"
                className="font-medium text-foreground underline underline-offset-4 hover:text-primary"
              >
                Sign up
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default SingIn;
