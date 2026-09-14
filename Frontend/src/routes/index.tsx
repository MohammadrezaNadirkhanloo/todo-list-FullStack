import SpinnerComponentPage from "@/components/SpinnerPage";
import { createBrowserRouter } from "react-router";

const lazyPage =
  (importer: () => Promise<{ default: React.ComponentType }>) => async () => ({
    Component: (await importer()).default,
  });

export const router = createBrowserRouter([
  {
    path: "/sign-in",
    // loader: authLoader,
    HydrateFallback: () => <SpinnerComponentPage />,
    // errorElement: <ErrorBoundary />,
    lazy: lazyPage(() => import("@/page/SingInPage")),
  },
  {
    path: "/sign-up",
    // loader: authLoader,
    HydrateFallback: () => <SpinnerComponentPage />,
    // errorElement: <ErrorBoundary />,
    lazy: lazyPage(() => import("@/page/SingUpPage")),
  },
  {
    path: "/",
    // loader: authLoader,
    HydrateFallback: () => <SpinnerComponentPage />,
    // errorElement: <ErrorBoundary />,
    lazy: lazyPage(() => import("@/page/Home")),
  },
  {
    path: "/todo-list",
    // loader: authLoader,
    HydrateFallback: () => <SpinnerComponentPage />,
    // errorElement: <ErrorBoundary />,
    lazy: lazyPage(() => import("@/page/TodoList")),
  },
  // {
  //   loader: protectedLoader,
  //   HydrateFallback: SpinnerComponentPage,
  //   children: [
  //     {
  //       path: "/",
  //       Component: RootLayout,
  //       children: [
  //         {
  //           index: true,
  //           loader: async () => {
  //             await useSessionStore.getState().ensureSession();
  //             const homeRoute = useSessionStore.getState().homeRoute;
  //             return redirect(homeRoute ?? "/dashboard");
  //           },
  //         },
  //         {
  //           path: "dashboard",
  //           loader: requireAbility("read", "Dashboard"),
  //           lazy: lazyPage(() => import("@/pages/Dashboard")),
  //         },
  //       ],
  //     },
  //   ],
  // },
  // {
  //   path: "403",
  //   lazy: lazyPage(() => import("@/pages/public/Forbidden")),
  // },
  {
    path: "*",
    lazy: lazyPage(() => import("@/page/NotFoundPage")),
  },
]);
