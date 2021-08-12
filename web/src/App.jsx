import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom'

import Header from './components/Header'

import '@/assets/styles/App.less'

import routes from '@/router'


function App() {


  return (
    <div className="App">
      <Router>
        <Header></Header>
        <Switch>
          {
            routes.map(route => <Route exact key={route.path} path={route.path}>
              <route.component />
            </Route>)
          }
        </Switch>
      </Router>
    </div>
  )
}

export default App
