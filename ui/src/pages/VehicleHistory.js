import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Moment from 'react-moment';

export default function VehicleHistory() {
  const { registrationNumber } = useParams();
  const [vehicle, setVehicle] = useState(null);

  useEffect(()=> {
    fetch('/vehicles', { 'method' : 'get' })
      .then(response => response.json())
      .then(vehicles => findVehicle(vehicles))
      .then(vehicle => setVehicle(vehicle))
  }, [registrationNumber]);

  function findVehicle(vehicles) {
    return vehicles.find(vehicle => vehicle.RegistrationNumber === registrationNumber);
  }

  function MOTs() {
    if(vehicle != null) {
      return vehicle.MOTHistory.map(mot => <MOTTest key={mot.TestNumber} {...mot}/>);
    } else {
      return null;
    }
  }

  return(
    <div className='container'>
      <div className='row title-row'>
        <div className='column'>
          <h1>{registrationNumber}</h1>
          <MakeAndModel vehicle={vehicle}/>
        </div>
        <div className='column'>
        </div>
      </div>

      {MOTs()}
    </div>
  );
}

function MakeAndModel({ vehicle }) {
  if(vehicle != null) {
    console.log(vehicle)
    return(<h3>{vehicle.Manufacturer} {vehicle.Model}</h3>)
  } else {
    return null;
  }
}

const commentTypes = {
  'FAIL': 'Reasons for failure',
  'DANGEROUS': 'Repair immediately (dangerous)',
  'MAJOR': 'Repair immediately (major)',
  'MINOR': 'Repair as soon as possible (minor)',
  'ADVISORY': 'Monitor and repair if necessary (advisory)',
  'PRS': 'Pass with Rectification',
  'USER ENTERED': 'Other comments'
};
const commentTypeOrder = Object.keys(commentTypes);

function MOTTest({ Passed, Mileage, ExpiryDate, CompletedDate, RfrAndComments }) {
  const commentsByType = (RfrAndComments || [])
    .sort((a, b) => commentTypeOrder.indexOf(a.Type) - commentTypeOrder.indexOf(b.Type))
    .reduce((accumulator, comment) => {
      if (accumulator[comment.Type] === undefined) {
        accumulator[comment.Type] = [];
      }
      accumulator[comment.Type].push(comment.Comment);

      return accumulator;
    }, {});

  const commentComponents = Object.entries(commentsByType)
    .map(([type, comments]) => {
      return(<CommentsList type={type} comments={comments} key={type} />);
    });

  return(
    <div className='row mot-test'>
      <div className='column pass-fail'>
        <PassOrFail Passed={Passed} />
      </div>

      <div className='column'>
        <div className='row'>
          <div className='column'>
            <label>Test Date</label>
            <Moment format='DD/MM/YYYY'>{CompletedDate}</Moment>
          </div>
          <div className='column'>
            <label>Mileage</label>
            {Mileage}
          </div>
          <div className='column'>
            <label>Expiry Date</label>
            <Moment format='DD/MM/YYYY'>{ExpiryDate}</Moment>
          </div>
        </div>

        { commentComponents }
      </div>
    </div>
  )
}

function PassOrFail({ Passed }) {
  if (Passed) {
    return(<h4 className='pass'>Pass</h4>);
  } else {
    return(<h4 className='fail'>Fail</h4>);
  }
}

function CommentsList({ type, comments }) {
  const commentComponents = comments.map((comment, i) => <Comment Comment={comment} key={`${type}-${i}`} />);

  const title = commentTypes[type] || 'Other comments';

  return(
    <div className='row comments-list'>
      <div className='column'>
        <label>{title}</label>
        <ul>
          {commentComponents}
        </ul>
      </div>
    </div>
  )
}

function Comment({ Comment }) {
  return(
    <li>{Comment}</li>
  );
}
