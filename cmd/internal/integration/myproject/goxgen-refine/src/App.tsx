import {GitHubBanner, Refine, ResourceProps, WelcomePage} from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";

import {Layout, List, notificationProvider, RefineSnackbarProvider, ThemedLayout, ThemedLayoutV2} from "@refinedev/mui";

import CssBaseline from "@mui/material/CssBaseline";
import GlobalStyles from "@mui/material/GlobalStyles";
import {gql, GraphQLClient} from "@refinedev/graphql";
import routerBindings, {
  DocumentTitleHandler,
  NavigateToResource,
  UnsavedChangesNotifier,
} from "@refinedev/react-router-v6";
import { useTranslation } from "react-i18next";
import {BrowserRouter, Outlet, Route, RouteProps, Routes} from "react-router-dom";
import { ColorModeContextProvider } from "./contexts/color-mode";
import {ReactElement, useState} from "react";
import Typography from "@mui/material/Typography";
import { OffLayoutArea } from "./components/off-layout-area";
import { GetDefaultActionConfig } from "./utils/xgen_graphql_defaults";
import {useXgenSet} from "./utils/xgen_resources_hook";
import xgenDataProvider from "./utils/xgen-data-provider";

const API_URL = "http://localhost:80/query";

const client = new GraphQLClient(API_URL);
const gqlDataProvider = xgenDataProvider(client);

function App() {
  const { t, i18n } = useTranslation();

  const i18nProvider = {
    translate: (key: string, params: object) => t(key, params),
    changeLocale: (lang: string) => i18n.changeLanguage(lang),
    getLocale: () => i18n.language,
  };


  const {resources, routes} = useXgenSet(client);

  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ColorModeContextProvider>
          <CssBaseline />
          <GlobalStyles styles={{ html: { WebkitFontSmoothing: "auto" } }} />
          <RefineSnackbarProvider>
            <Refine
              dataProvider={gqlDataProvider}
              notificationProvider={notificationProvider}
              routerProvider={routerBindings}
              i18nProvider={i18nProvider}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
              }}
              resources={resources}
            >
              <Routes>
                <Route
                    element={
                      <ThemedLayoutV2
                          OffLayoutArea={OffLayoutArea}
                      >
                        <Outlet />
                      </ThemedLayoutV2>
                    }
                >
                  {/*<Route index element={<NavigateToResource resource="blog_posts" />} />*/}
                  {routes.map((route, index) => (
                      <Route key={index} {...route} />
                  ))}
                  {/*<Route path="user">*/}
                  {/*  <Route index element={<div />} />*/}
                  {/*  <Route*/}
                  {/*      path="show/:id"*/}
                  {/*      element={<div />}*/}
                  {/*  />*/}
                  {/*  <Route*/}
                  {/*      path="edit/:id"*/}
                  {/*      element={<div />}*/}
                  {/*  />*/}
                  {/*  <Route*/}
                  {/*      path="create"*/}
                  {/*      element={<div />}*/}
                  {/*  />*/}
                  {/*</Route>*/}
                </Route>
              </Routes>
              <RefineKbar />
              <UnsavedChangesNotifier />
              <DocumentTitleHandler />
            </Refine>
          </RefineSnackbarProvider>
        </ColorModeContextProvider>
      </RefineKbarProvider>
    </BrowserRouter>
  );
}

export default App;
