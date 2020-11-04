<?php

namespace Tests\Unit\Models;

use App\Models\AbstractDoctrineRepository;
use App\Models\DatabaseHandler;
use Doctrine\Common\Collections\ArrayCollection;
use Mockery as m;
use Mockery\MockInterface;
use Tests\Helper\ModelHelper;
use Tests\TestCase;

/**
 * Class AbstractDoctrineRepositoryTest
 *
 * @package Tests\Unit\Models
 */
final class AbstractDoctrineRepositoryTest extends TestCase
{
    use ModelHelper;

    //region Tests

    private function setUpFindTest(bool $withModel = true): array
    {
        $id = $this->getFaker()->numberBetween();
        $model = $this->createModel();
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerFind($databaseHandler, $withModel ? $model : null, $id);
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $id, $model];
    }

    /**
     * @return void
     */
    public function testFind(): void
    {
        [$doctrineRepository, $id, $model] = $this->setUpFindTest();

        $this->assertEquals($model, $doctrineRepository->find($id));
    }

    /**
     * @return void
     */
    public function testFindWithoutModel(): void
    {
        [$doctrineRepository, $id] = $this->setUpFindTest(false);

        $this->assertEmpty($doctrineRepository->find($id));
    }

    private function setUpFindAllTest(bool $withModel = true): array
    {
        $models = new ArrayCollection([$this->createModel()]);
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerLoadAll($databaseHandler, $withModel ? $models : new ArrayCollection());
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $models];
    }

    /**
     * @return void
     */
    public function testFindAll(): void
    {
        [$doctrineRepository, $models] = $this->setUpFindAllTest();

        $this->assertEquals($models, $doctrineRepository->findAll());
    }

    /**
     * @return void
     */
    public function testFindAllWithoutModels(): void
    {
        [$doctrineRepository] = $this->setUpFindAllTest(false);

        $this->assertEquals(new ArrayCollection(), $doctrineRepository->findAll());
    }

    private function setUpFindOneByTest(bool $withModel = true): array
    {
        $criteria = [$this->getFaker()->word => $this->getFaker()->word];
        $model = $this->createModel();
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerLoad($databaseHandler, $withModel ? $model : null, $criteria);
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $criteria, $model];
    }

    /**
     * @return void
     */
    public function testFindOneBy(): void
    {
        [$doctrineRepository, $criteria, $model] = $this->setUpFindOneByTest();

        $this->assertEquals($model, $doctrineRepository->findOneBy($criteria));
    }

    /**
     * @return void
     */
    public function testFindOneByWithoutModel(): void
    {
        [$doctrineRepository, $criteria] = $this->setUpFindOneByTest(false);

        $this->assertNull($doctrineRepository->findOneBy($criteria));
    }


    private function setUpFindByTest(bool $withModels = true): array
    {
        $criteria = [$this->getFaker()->word => $this->getFaker()->word];
        $orderBy = [$this->getFaker()->word => $this->getFaker()->word];
        $limit = $this->getFaker()->numberBetween();
        $offset = $this->getFaker()->numberBetween();
        $models = new ArrayCollection($withModels ? [$this->createModel()] : []);
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerLoadAll($databaseHandler, $models, $criteria, $orderBy, $limit, $offset);
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $criteria, $orderBy, $limit, $offset, $models];
    }

    /**
     * @return void
     */
    public function testFindBy(): void
    {
        [$doctrineRepository, $criteria, $orderBy, $limit, $offset, $models] = $this->setUpFindByTest();

        $this->assertEquals($models, $doctrineRepository->findBy($criteria, $orderBy, $limit, $offset));
    }

    /**
     * @return void
     */
    public function testFindByWithoutModels(): void
    {
        [$doctrineRepository, $criteria, $orderBy, $limit, $offset] = $this->setUpFindByTest(false);

        $this->assertEmpty($doctrineRepository->findBy($criteria, $orderBy, $limit, $offset));
    }

    private function setUpSaveTest(bool $flush = null): array
    {
        $model = $this->createModel();
        $databaseHandler = $this->createDatabaseHandler();
        $this->mockDatabaseHandlerSave($databaseHandler, $model, $flush === null ? true : $flush);
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $model, $databaseHandler];
    }

    /**
     * @return void
     */
    public function testSave(): void
    {
        [$doctrineRepository, $model, $databaseHandler] = $this->setUpSaveTest();

        $this->assertEquals($model, $doctrineRepository->save($model));
        $this->assertDatabaseHandlerSave($databaseHandler, $model, true);
    }

    /**
     * @return void
     */
    public function testSaveWithFlush(): void
    {
        $flush = $this->getFaker()->boolean;
        [$doctrineRepository, $model, $databaseHandler] = $this->setUpSaveTest($flush);

        $this->assertEquals($model, $doctrineRepository->save($model, $flush));
        $this->assertDatabaseHandlerSave($databaseHandler, $model, $flush);
    }

    private function setUpDeleteTest(): array
    {
        $model = $this->createModel();
        $databaseHandler = $this->createDatabaseHandler();
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        return [$doctrineRepository, $model, $databaseHandler];
    }

    /**
     * @return void
     */
    public function testDelete(): void
    {
        [$doctrineRepository, $model, $databaseHandler] = $this->setUpDeleteTest();

        $this->assertEquals($doctrineRepository, $doctrineRepository->delete($model));
        $this->assertDatabaseHandlerDelete($databaseHandler, $model, true);
    }

    /**
     * @return void
     */
    public function testDeleteWithFlush(): void
    {
        $flush = $this->getFaker()->boolean;
        [$doctrineRepository, $model, $databaseHandler] = $this->setUpDeleteTest();

        $doctrineRepository->delete($model, $flush);

        $this->assertDatabaseHandlerDelete($databaseHandler, $model, $flush);
    }

    /**
     * @return void
     */
    public function testFlush(): void
    {
        $databaseHandler = $this->createDatabaseHandler();
        $doctrineRepository = $this->getAbstractDoctrineRepository($databaseHandler);

        $this->assertEquals($doctrineRepository, $doctrineRepository->flush());
        $this->assertDatabaseHandlerFlush($databaseHandler);
    }

    //endregion

    /**
     * @param DatabaseHandler|null $databaseHandler
     *
     * @return AbstractDoctrineRepository|MockInterface
     */
    private function getAbstractDoctrineRepository(DatabaseHandler $databaseHandler = null): AbstractDoctrineRepository
    {
        return m::mock(AbstractDoctrineRepository::class, [$databaseHandler ?: $this->createDatabaseHandler()])
            ->makePartial()
            ->shouldAllowMockingProtectedMethods();
    }
}
