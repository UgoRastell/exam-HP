func (s *server) StreamSensors(req *pb.Empty, stream pb.SensorService_StreamSensorsServer) error {
    for data := range dataChan {
        if err := stream.Send(&pb.SensorData{
            Id:          data.ID,
            Temperature: data.Temperature,
            Humidity:    data.Humidity,
            Timestamp:   data.Timestamp,
        }); err != nil {
            return err
        }
    }
    return nil
}
