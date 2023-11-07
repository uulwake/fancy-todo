export class CustomError extends Error {
  status: number;

  constructor(data: { message: string; status: number }) {
    super(data.message);
    this.status = data.status;
  }
}
