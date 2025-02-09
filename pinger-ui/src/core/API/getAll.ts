import axios from "axios";
import {FetchResponse, IContainer} from "../../types/types.ts";
const GetAll = async (): Promise<IContainer[]> => {
    const response = await axios.get<FetchResponse>("http://localhost:8080/api/v1/getAll");
    return response.data.containers; // Возвращаем только массив контейнеров
};

export default GetAll
