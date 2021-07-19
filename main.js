//CRUD
const returnAllTickets = document.querySelector(".tickets tbody");
const createTicket = document.querySelector("#create-ticket");
const updateTicket = document.querySelector("#edit-ticket-modal");
const deleteTicket = document.querySelector("#delete-ticket-modal");

//TICKET
const formIssue = document.getElementById("form-issue");
const formPriority = document.getElementById("form-priority");
const formStatus = document.getElementById("form-status");

//TICKET MODAL
const modalIssue = document.getElementById("modal-issue");
const modalPriority = document.getElementById("modal-priority");
const modalStatus = document.getElementById("modal-status");


const renderTicket = (ticket) => {
  let output = "";
  ticket.forEach((ticket) => {
    output += `
        <tr>
            <td id="ticketId">${ticket.id}</td>
            <td id="ticketIssue">${ticket.issue}</td>
            <td id="ticketPriority">${ticket.priority}</td>
            <td id="ticketStatus">${ticket.status}</td>
            <td>
                <a href="#edit-ticket-modal" class="edit" data-toggle="modal"><i class="material-icons" data-toggle="tooltip" title="Edit">&#xE254;</i></a>
                <a href="#delete-ticket-modal" class="delete" data-toggle="modal"><i class="material-icons" data-toggle="tooltip" title="Delete">&#xE872;</i></a>
            </td>
        </tr>
        `;
  });
  returnAllTickets.innerHTML = output;
};

const url = "http://localhost:8080";

//RETURN ALL TICKETS
fetch(url, {
  method: "GET",
})
  .then((res) => res.json())
  .then((data) => renderTicket(data));

//CREATE TICKET
createTicket.addEventListener("submit", (e) => {
  e.preventDefault();

  fetch(`${url}/ticket`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      issue: formIssue.value,
      priority: formPriority.value,
      status: formStatus.value,
    }),
  }).then(() => location.reload());
});

//DELETE/EDIT TICKET
returnAllTickets.addEventListener("click", (e) => {
  e.preventDefault();
  const id = e.target.closest("tr").querySelector("td#ticketId").textContent;
  let issue = e.target.closest("tr").querySelector("td#ticketIssue").textContent;

  let deleteTicketIsPressed = e.target.title == "Delete";
  let editTickerIsPressed = e.target.title == "Edit";
  
  //DELETE TICKET
  if (deleteTicketIsPressed) {
    deleteTicket.addEventListener("submit", (e) => {
      fetch(`${url}/ticket/${id}`, {
        method: "DELETE",
      }).then(() => location.reload());
    });
  }
  
  //EDIT TICKET
  if (editTickerIsPressed) {
    modalIssue.value = issue;

    updateTicket.addEventListener("submit", (e) => {
      e.preventDefault();


      fetch(`${url}/ticket/update`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: parseInt(id),
          issue: modalIssue.value,
          priority: modalPriority.value,
          status: modalStatus.value,
        }),
      }).then(() => location.reload());
    });
  }
});
