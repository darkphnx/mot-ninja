import { useEffect, useState } from 'react';
import Moment from 'react-moment';
import moment from 'moment';
import { Link } from "react-router-dom";

export default function VehicleList() {
  const [addFormVisible, setAddFormVisible] = useState(true);

  return(
    <div className="container">
      <div className='row title-row'>
        <div className='column'>
          <h1>Your Vehicles</h1>
        </div>
        <div className='column'>
          <AddVehicleForm/>
        </div>
      </div>

      <div class='row'>
        <VehicleTable />
      </div>
    </div>
  )
}

function AddVehicleForm() {
  return (
    <div className='row add-vehicle'>
      <div className='column'>
        <input type='text' id='registration-number' placeholder='Registration Number' />
      </div>

      <div className='column add-vehicle-submit'>
        <a href='#' className='button'>Add Vehicle</a>
      </div>
    </div>
  )
}

function VehicleTable() {
  const [vehicles, setVehicles] = useState([]);

  useEffect(()=> {
    fetch('/vehicles', { 'method' : 'get' })
      .then(response => response.json())
      .then(vehicles => setVehicles(vehicles))
  }, []);

  function vehicleComponents() {
    return vehicles.map((vehicle, i) => {
      return(<Vehicle key={vehicle.ID} {...vehicle} />);
    })
  }

  return(
    <table>
      <thead>
        <tr>
          <th>Registration</th>
          <th>Make/Model</th>
          <th>MOT Status</th>
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
    if(latestMOT != null) {
      const advisories = latestMOT.RfrAndComments.filter(comment => comment.Type === 'MINOR')

      return advisories.length;
    }
  }


  return(
    <tr>
      <td>
        <Link to={"/" + ID}>{RegistrationNumber}</Link>
      </td>
      <td>{Manufacturer} {Model}</td>
      <td>{expiredOrDue(MotDue)} <Moment format='DD/MM/YYYY'>{MotDue}</Moment></td>
      <td>{advisoryCount()}</td>
    </tr>
  );
}
