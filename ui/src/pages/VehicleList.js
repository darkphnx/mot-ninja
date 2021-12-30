import { useEffect, useState } from 'react';
import Moment from 'react-moment';
import moment from 'moment';
import { Link } from "react-router-dom";
import FormErrors from '../components/FormErrors';

export default function VehicleList() {
  const [vehicles, setVehicles] = useState([]);
  const [searchFilter, setSearchFilter] = useState("");

  useEffect(()=> {
    fetch('/api/vehicles', { method: 'GET' })
      .then(response => response.json())
      .then(vehicles => setVehicles(vehicles || []));
  }, []);

  function handleOnVehicleAdded(addedVehicle) {
    setVehicles([...vehicles, addedVehicle]);
  }

  function handleSearchQueryUpdated(newSearchQuery) {
    setSearchFilter(newSearchQuery);
  }

  // TODO: Replace with a proper database search
  function filteredVehicles() {
    if(searchFilter === "") {
      return vehicles;
    } else {
      const filter = searchFilter.toLowerCase();
      const searchFields = ['RegistrationNumber', 'Manufacturer', 'Model'];

      return vehicles.filter(vehicle => {
        return searchFields.some(field => {
          return vehicle[field].toLowerCase().startsWith(filter);
        });
      });
    }
  }

  return(
    <div className="container">
      <div className='row title-row'>
        <div className='column'>
          <h1>Your Vehicles</h1>
        </div>
        <div className='column'>
          <SearchVehicleForm onSearchUpdate={handleSearchQueryUpdated} />
        </div>
      </div>

      <div className='row vehicle-list-content'>
        <div className='column'>
          <VehicleTable vehicles={filteredVehicles()}/>
        </div>

        <div className='column add-vehicle'>
          <h4>Add a new vehicle</h4>
          <AddVehicleForm onVehicleAdded={handleOnVehicleAdded} />
          <p>You may monitor up to five vehicles.</p>
        </div>
      </div>
    </div>
  )
}

function SearchVehicleForm({ onSearchUpdate }) {
  const [searchQuery, setSearchQuery] = useState("")

  function handleSearchQuery(e) {
    setSearchQuery(e.target.value);
    onSearchUpdate(e.target.value);
  }

  return (
    <div className='row search-vehicle'>
      <div className='column'>
        <input type='text' id='query' placeholder='Search Registration/Make/Model' value={searchQuery} onChange={handleSearchQuery}/>
      </div>
    </div>
  );
}

function AddVehicleForm({ onVehicleAdded }) {
  const [registrationNumber, setRegistrationNumber] = useState("")
  const [formErrors, setFormErrors] = useState([])

  function handleRegistrationNumber(e) {
    setRegistrationNumber(e.target.value);
  }

  function submitForm(e) {
    fetch('/api/vehicles', {
      method: 'POST',
      body: JSON.stringify({ "RegistrationNumber" : registrationNumber })
    }).then(response => response.json())
      .then(payload => {
        if (payload.Error) {
          setFormErrors(payload.Error)
        } else {
          handleVehicleAdded(payload)
        }
      });
  }

  function handleVehicleAdded(vehicle) {
    onVehicleAdded(vehicle);
    setRegistrationNumber("");
  }

  return (
    <div>
      <FormErrors errors={formErrors} />
      <div className='add-vehicle-form row'>
        <div className='column'>
          <input type='text' id='registration-number' placeholder='Registration Number' value={registrationNumber} onChange={handleRegistrationNumber}/>
        </div>

        <div className='column'>
          <a href='#' className='button' onClick={submitForm}>Add Vehicle</a>
        </div>
      </div>
    </div>
  )
}

function VehicleTable({ vehicles }) {
  function vehicleComponents() {
    return (vehicles).map((vehicle, i) => {
      return(<Vehicle key={vehicle.ID} {...vehicle} />);
    })
  }

  return(
    <table>
      <thead>
        <tr>
          <th>Registration</th>
          <th>Make/Model</th>
          <th>MOT Due</th>
          <th>Advisories</th>
        </tr>
      </thead>
      <tbody>
        {vehicleComponents()}
      </tbody>
    </table>
  )
}

function Vehicle({ ID, RegistrationNumber, Manufacturer, Model, MotDue, VEDDue, MOTHistory }) {
  function expiredOrDue(timestamp) {
    if(moment(timestamp).isBefore(moment())){
      return 'Expired';
    } else {
      return 'Due';
    }
  }

  function findLatestMOT() {
    if(MOTHistory == null) {
      return null;
    }

    const sortedHistory = MOTHistory.sort((a, b) => {
      if(a.CompletedDate > b.CompletedDate) {
        return -1;
      } else if (a.CompletedDate < b.CompletedDate) {
        return 1;
      } else {
        return 0;
      }
    });

    return sortedHistory[0];
  }

  const latestMOT = findLatestMOT();

  function advisoryCount() {
    if(latestMOT != null && latestMOT.RfrAndComments != null) {
      const advisories = latestMOT.RfrAndComments.filter(comment => comment.Type === 'MINOR')

      return advisories.length;
    } else {
      return 0;
    }
  }


  return(
    <tr>
      <td>
        <Link to={"/" + RegistrationNumber}>{RegistrationNumber}</Link>
      </td>
      <td>{Manufacturer} {Model}</td>
      <td>{expiredOrDue(MotDue)} <Moment format='DD/MM/YYYY'>{MotDue}</Moment></td>
      <td>{advisoryCount()}</td>
    </tr>
  );
}
