import './App.css';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";
import VehicleList from './components/VehicleList';
import VehicleHistory from './components/VehicleHistory';

function App() {
  return (
    <Router>
      <div className="container">
        <div className="row">
          <h1>MOT Minder</h1>
        </div>
        <div className="row">
          <Switch>
            <Route path="/:id">
              <VehicleHistory />
            </Route>
            <Route path="/">
              <VehicleList />
            </Route>
          </Switch>
        </div>
      </div>
    </Router>
  );
}

export default App;
