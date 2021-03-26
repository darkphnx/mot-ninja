import './App.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Signup from './pages/Signup';
import VehicleList from './pages/VehicleList';
import VehicleHistory from './pages/VehicleHistory';

import logo from './images/logo.svg'

function App() {
  return (
    <Router>
      <div className="wrapper">
        <nav className="navigation">
          <div className="container">
            <Link to="/" className="navigation-title">
              <img src={logo} alt="logo" className="navigation-logo" />
              MOT.ninja
            </Link>
          </div>
        </nav>

        <main className="content">
          <Switch>
            <Route path="/signup">
              <Signup />
            </Route>
            <Route path="/:registrationNumber">
              <VehicleHistory />
            </Route>
            <Route path="/">
              <VehicleList />
            </Route>
          </Switch>
        </main>
      </div>
    </Router>
  );
}

export default App;
