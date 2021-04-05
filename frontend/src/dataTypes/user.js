export class UserPointer {
  constructor(id = null, firstName = null, lastName = null) {
    this.id = id;
    (this.firstName = firstName), (this.lastName = lastName);
  }
}
