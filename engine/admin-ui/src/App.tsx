import 'Styles/app.global.scss';
import 'kwc/index.css';
import 'markdown-navbar/dist/navbar.css';
import 'Styles/markdown-navbar.scss';
import 'react-toastify/dist/ReactToastify.css';
import 'Styles/react-toastify.scss';

import { ErrorMessage, SpinnerCircular } from 'kwc';
import { Redirect, Switch } from 'react-router';
import { Route, Router } from 'react-router-dom';

import AccessDenied from 'Pages/AccessDenied/AccessDenied';
import AddApiToken from 'Pages/AddApiToken/AddApiToken';
import AddUser from 'Pages/AddUser/AddUser';
import AddVersion from 'Pages/AddVersion/AddVersion';
import { GetMe } from 'Graphql/queries/types/GetMe';
import { GetVersionConfStatus } from 'Graphql/queries/types/GetVersionConfStatus';
import { GlobalHotKeys } from 'react-hotkeys';
import Login from 'Pages/Login/Login';
import Logs from 'Pages/Logs/Logs';
import LogsPanel from 'Pages/Version/pages/Status/LogsPanel/LogsPanel';
import MagicLink from 'Pages/MagicLink/MagicLink';
import NotFound from 'Pages/NotFound/NotFound';
import NotificationService from 'Components/NotificationService/NotificationService';
import Profile from 'Pages/Profile/Profile';
import ROUTE from 'Constants/routes';
import React from 'react';
import Runtime from 'Pages/Runtime/Runtime';
import Settings from 'Pages/Settings/Settings';
import SubscriptionService from 'Components/SubscriptionService/SubscriptionService';
import Users from 'Pages/Users/Users';
import UsersActivity from 'Pages/UsersActivity/UsersActivity';
import VerifyEmail from 'Pages/VerifyEmail/VerifyEmail';
import { buildRoute } from 'Utils/routes';
import { getNotAllowedRoutes } from './accessLevelRoutes';
import history from './browserHistory';
import keymaps from './keymaps';
import { loader } from 'graphql.macro';
import useLogin from 'Graphql/hooks/useLogin';
import useLogs from 'Graphql/hooks/useLogs';
import { useQuery } from '@apollo/client';

const GetMeQuery = loader('Graphql/queries/getMe.graphql');
const GetRuntimeQuery = loader('Graphql/queries/getRuntimeAndVersions.graphql');

function ProtectedRoutes() {
  const { data, error, loading } = useQuery<GetMe>(GetMeQuery);
  const { data: runtimeData, loading: runtimeLoading } = useQuery<
    GetVersionConfStatus
  >(GetRuntimeQuery);
  const { login } = useLogin();

  if (loading || runtimeLoading || !runtimeData) {
    return (
      <div className="splash">
        <SpinnerCircular />
      </div>
    );
  }

  if (error || !data || !data.me || !runtimeData.runtime) {
    return <ErrorMessage />;
  }

  login();

  const protectedRoutes: string[] = getNotAllowedRoutes(data.me.accessLevel);

  return (
    <>
      <div className="page-with-logs-wrapper">
        <Switch>
          <Route path={protectedRoutes} component={AccessDenied} />
          <Redirect exact from={ROUTE.SETTINGS} to={ROUTE.SETTINGS_GENERAL} />
          <Redirect
            exact
            from={ROUTE.RUNTIME_VERSION}
            to={ROUTE.RUNTIME_VERSION_STATUS}
          />
          <Redirect exact from={ROUTE.RUNTIME} to={ROUTE.RUNTIME_VERSIONS} />

          <Route path={ROUTE.NEW_USER} component={AddUser} />
          <Route path={ROUTE.NEW_API_TOKEN} component={AddApiToken} />
          <Route path={ROUTE.NEW_VERSION} component={AddVersion} />

          <Redirect
            exact
            from={ROUTE.HOME}
            to={buildRoute.runtime(ROUTE.RUNTIME, runtimeData.runtime.id)}
          />

          <Route exact path={ROUTE.LOGS} component={Logs} />
          <Route
            path={[
              ROUTE.RUNTIME_VERSION_CONFIGURATION,
              ROUTE.RUNTIME_VERSION_METRICS,
              ROUTE.RUNTIME_VERSION_DOCUMENTATION,
              ROUTE.RUNTIME_VERSION_STATUS
            ]}
            component={Runtime}
          />
          <Route path={ROUTE.RUNTIME_VERSIONS} component={Runtime} />

          <Route path={ROUTE.SETTINGS} component={Settings} />
          <Route path={ROUTE.PROFILE} component={Profile} />
          <Route path={ROUTE.USERS} component={Users} />
          <Route path={ROUTE.AUDIT} component={UsersActivity} />
          <Route component={NotFound} />
        </Switch>
        <LogsPanel />
      </div>
    </>
  );
}

export function Routes() {
  return (
    <Switch>
      <Route exact path={ROUTE.LOGIN} component={Login} />
      <Route exact path={ROUTE.VERIFY_EMAIL} component={VerifyEmail} />
      <Route exact path={ROUTE.MAGIC_LINK} component={MagicLink} />

      <Route component={ProtectedRoutes} />
    </Switch>
  );
}

function App() {
  const { openLogs, closeLogs } = useLogs();

  const handlers = {
    CLOSE_LOGS: closeLogs,
    OPEN_LOGS: openLogs
  };

  return (
    <GlobalHotKeys keyMap={keymaps} handlers={handlers}>
      <div className="app">
        <Router history={history}>
          <SubscriptionService />
          <NotificationService />
          <Routes />
        </Router>
      </div>
      <div id="chartjs-tooltip">
        <table />
      </div>
    </GlobalHotKeys>
  );
}

export default App;
