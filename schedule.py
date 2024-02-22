from math import ceil, floor

BATCH_SIZE = 100


class Batch:
    def __init__(self, dataset_len, idx=0, device_count=1):
        self.__batch_size = BATCH_SIZE
        self.__device_count = device_count
        self.__dataset_len = dataset_len
        self.__idx = {"f":idx,"l":idx}

        self.__init_partition_value
        self.__server_ratio
        self.__device_ratio

        self.__Set_ratio(self.__batch_size)

    def __Set_ratio(self, size):
        self.__init_partition_value = size / (self.__device_count + 1)

        self.__server_ratio = ceil(self.__init_partition_value)
        self.__device_ratio = [
            floor(self.__init_partition_value) for _ in range(self.__device_count)
        ]

    def __Batch(self):
        if self.__idx + self.__batch_size - 1 > self.__dataset_len:
            self.__idx = [self.__idx["l"], self.__dataset_len]
        else:
            self.__idx = [self.__idx["l"], self.__idx["l"] + self.__batch_size]

    def __Set_size(self):
        if self.__batch_size

        
    def Batch_partition(self):
        self.__Batch()
        
        if  
        server_idx = None
        device_idx = None

        return server_idx, device_idx
