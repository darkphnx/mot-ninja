import './App.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import VehicleList from './components/VehicleList';
import VehicleHistory from './components/VehicleHistory';

function App() {
  return (
    <Router>
      <div className="wrapper">
        <nav className="navigation">
          <div className="container">
            <Link to="/" className="navigation-title">MOT Minder</Link>
          </div>
        </nav>

        <main className="content">
          <Switch>
            <Route path="/:id">
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
