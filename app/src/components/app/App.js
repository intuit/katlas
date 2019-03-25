import React, { Component } from "react";
import { Route, Switch } from "react-router-dom";

import MenuBar from "../menuBar/MenuBar";
import Home from "../home/Home";
import Results from "../results/Results";
import GraphContainer from "../graph/GraphContainer";
import Notifier from '../notifier/Notifier';
import ErrorPage from '../errorPage/ErrorPage';

export default class App extends Component {
  render() {
    return (
      <div className="App">
        <MenuBar/>
        <Notifier/>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/results" component={Results} />
          <Route path="/graph/:uidOrQsl" component={GraphContainer} />
          {/* catch-all route in case of no earlier match */}
          <Route component={ErrorPage} />
        </Switch>
      </div>
    );
  }
}