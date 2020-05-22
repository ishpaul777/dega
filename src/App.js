import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import 'antd/dist/antd.css';
import BasicLayout from './layouts/basic';

//Pages
import Spaces from './pages/spaces';
import Media from './pages/media';
import Dashboard from './pages/dashboard';
import CreateSpace from './pages/spaces/create';
function App() {
  return (
    <div className="App">
      <Router basename={process.env.PUBLIC_URL}>
        <BasicLayout>
          <Switch>
            <Route exact path="/" component={Dashboard} />
            <Route path="/dashboard" component={Dashboard} />
            <Route exact path="/spaces" component={Spaces} />
            <Route exact path="/media" component={Media} />
            <Route path="/spaces/create" component={CreateSpace} />
          </Switch>
        </BasicLayout>
      </Router>
    </div>
  );
}

export default App;
