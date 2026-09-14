import { cn } from "@/lib/utils";

import { GithubIcon } from "@/components/icons/github-icon";
import { GoogleIcon } from "@/components/icons/google-icon";
import { Button } from "@/components/ui/button";
import { AuthDivider } from "@/components/auth-divider";
import { DecorIcon } from "@/components/decor-icon";

import { SingUpForm } from "./form/FormLogin";
import { Link, useNavigate } from "react-router";
import ButtonComponent from "@/components/Button-custom";
import { HomeIcon } from "lucide-react";

function SingUp() {
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
                icon={<HomeIcon />}
                size={"icon"}
                className="mt-0"
                onClick={() => navigate("/")}
              />
              <h1 className="text-2xl font-bold tracking-wide">
                Create your account
              </h1>
            </div>

            <p className="text-base text-muted-foreground">
              Sign up to organize your tasks and get things done.
            </p>
          </div>

          {/* Sign Up Form */}
          <div className="space-y-4">
            <SingUpForm />

            <AuthDivider>OR</AuthDivider>

            {/* Social Login */}
            <div className="grid grid-cols-2 gap-2">
              <Button className="w-full" type="button" variant="outline">
                <GoogleIcon data-icon="inline-start" />
                Google
              </Button>

              <Button className="w-full" type="button" variant="outline">
                <GithubIcon data-icon="inline-start" />
                GitHub
              </Button>
            </div>
          </div>

          {/* Sign In */}
          <div className="space-y-4">
            <p className="text-center text-sm text-muted-foreground">
              Already have an account?{" "}
              <Link
                to="/sign-in"
                className="font-medium text-foreground underline underline-offset-4 hover:text-primary"
              >
                Sign in
              </Link>
            </p>

          
          </div>
        </div>
      </div>
    </div>
  );
}

export default SingUp;
