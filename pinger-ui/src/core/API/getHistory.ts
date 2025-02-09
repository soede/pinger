import axios from "axios";
import { FetchResponse, IContainer } from "../../types/types.ts";

const GetHistory = async (): Promise<IContainer[]> => {
    const response = await axios.get<FetchResponse>("http://localhost:8080/api/v1/history");
    return response.data.containers;
};

export default GetHistory;
